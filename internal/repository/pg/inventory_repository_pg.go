package pg

import (
	"context"
	"fmt"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryRepositoryPg struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewInventoryRepositoryPg(db *gorm.DB, log *logrus.Logger) repository.InventoryRepository {
	return &InventoryRepositoryPg{db: db, log: log}
}

func (r *InventoryRepositoryPg) GetStock(ctx context.Context, variantID string) (*entity.ProductVariantStock, error) {
	var stock entity.ProductVariantStock
	err := r.db.WithContext(ctx).Where("product_variant_id = ?", variantID).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *InventoryRepositoryPg) UpdateStock(ctx context.Context, variantID string, qtyChange int, reason string, orderID *string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stock entity.ProductVariantStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_variant_id = ?", variantID).
			First(&stock).Error; err != nil {
			return fmt.Errorf("failed to lock stock: %w", err)
		}

		newStock := stock.CurrentStock + qtyChange
		if newStock < 0 {
			return fmt.Errorf("insufficient stock: current %d, change %d", stock.CurrentStock, qtyChange)
		}

		if err := tx.Model(&stock).Update("current_stock", newStock).Error; err != nil {
			return fmt.Errorf("failed to update stock: %w", err)
		}

		ledger := entity.InventoryLedger{
			ProductVariantID: variantID,
			QuantityChange:   qtyChange,
			Reason:           reason,
			OrderID:          orderID,
		}
		if err := tx.Create(&ledger).Error; err != nil {
			return fmt.Errorf("failed to create ledger entry: %w", err)
		}

		return nil
	})
}
