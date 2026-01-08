package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type ProductVariantStockRepository interface {
	CreateStock(ctx context.Context, stock *entity.ProductVariantStock) error
	UpdateStock(ctx context.Context, stock *entity.ProductVariantStock, variantID string) error
	GetStockByVariantID(ctx context.Context, variantID string) (*entity.ProductVariantStock, error)
	// DeductStockWithVersion atomically decrements current_stock and reserved_stock with optimistic locking
	// Returns error if version mismatch (concurrent modification detected)
	DeductStockWithVersion(ctx context.Context, variantID string, quantity int, expectedVersion int) error
}
