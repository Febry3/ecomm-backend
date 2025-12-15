package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type ProductVariantStockRepository interface {
	CreateStock(ctx context.Context, stock *entity.ProductVariantStock) error
	UpdateStock(ctx context.Context, stock *entity.ProductVariantStock, variantID string) error
	GetStockByVariantID(ctx context.Context, variantID string) (*entity.ProductVariantStock, error)
}
