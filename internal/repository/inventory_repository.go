package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type InventoryRepository interface {
	GetStock(ctx context.Context, variantID string) (*entity.ProductVariantStock, error)
	UpdateStock(ctx context.Context, variantID string, qtyChange int, reason string, orderID *string) error
}
