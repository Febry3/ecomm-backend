package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type ProductVariantRepository interface {
	CreateProductVariant(ctx context.Context, productVariant *entity.ProductVariant) error
	DeleteProductVariant(ctx context.Context, productVariantID string) error
	GetProductVariant(ctx context.Context, productVariantID string) (*entity.ProductVariant, error)
	GetProductVariants(ctx context.Context) ([]entity.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, productVariant *entity.ProductVariant, productVariantID string) error
}
