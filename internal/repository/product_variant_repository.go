package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type ProductVariantRepository interface {
	CreateProductVariant(ctx context.Context, productVariant *entity.ProductVariant) error
	DeleteProductVariant(ctx context.Context, productVariantID string, sellerID int64) error
	GetProductVariant(ctx context.Context, productVariantID string) (*entity.ProductVariant, error)
	GetProductVariants(ctx context.Context, productID string) ([]entity.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, productVariant *entity.ProductVariant, productVariantID string) error
	GetProductVariantByID(ctx context.Context, productVariantID string) (*entity.ProductVariant, error)
}
