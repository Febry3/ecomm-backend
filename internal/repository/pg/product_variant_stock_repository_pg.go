package pg

import (
	"context"
	"errors"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

var ErrOptimisticLockFailed = errors.New("optimistic lock failed: version mismatch")

type ProductVariantStockRepositoryPg struct {
	db *gorm.DB
}

func NewProductVariantStockRepositoryPg(db *gorm.DB) repository.ProductVariantStockRepository {
	return &ProductVariantStockRepositoryPg{db: db}
}

func (r *ProductVariantStockRepositoryPg) CreateStock(ctx context.Context, stock *entity.ProductVariantStock) error {
	db := TxFromContext(ctx, r.db)
	return db.Create(stock).Error
}

func (r *ProductVariantStockRepositoryPg) UpdateStock(ctx context.Context, stock *entity.ProductVariantStock, variantID string) error {
	db := TxFromContext(ctx, r.db)
	return db.Model(&entity.ProductVariantStock{}).Where("product_variant_id = ?", variantID).Updates(stock).Error
}

func (r *ProductVariantStockRepositoryPg) GetStockByVariantID(ctx context.Context, variantID string) (*entity.ProductVariantStock, error) {
	var stock entity.ProductVariantStock
	err := r.db.WithContext(ctx).Where("product_variant_id = ?", variantID).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

// DeductStockWithVersion atomically decrements current_stock and reserved_stock with optimistic locking.
// This is called when payment is successful to convert reserved stock to actual stock deduction.
// Returns ErrOptimisticLockFailed if version mismatch (concurrent modification detected).
func (r *ProductVariantStockRepositoryPg) DeductStockWithVersion(ctx context.Context, variantID string, quantity int, expectedVersion int) error {
	db := TxFromContext(ctx, r.db)

	// Use raw SQL for atomic update with version check
	result := db.Exec(`
		UPDATE product_variant_stocks 
		SET current_stock = current_stock - ?, 
			reserved_stock = reserved_stock - ?,
			version = version + 1,
			last_updated = NOW()
		WHERE product_variant_id = ? AND version = ?
	`, quantity, quantity, variantID, expectedVersion)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrOptimisticLockFailed
	}

	return nil
}
