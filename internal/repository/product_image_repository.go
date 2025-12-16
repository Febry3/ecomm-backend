package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type ProductImageRepository interface {
	CreateProductImage(ctx context.Context, productImage *entity.ProductImage) error
}
