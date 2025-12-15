package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

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
