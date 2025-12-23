package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type ProductRepository interface {
	GetProductForBuyer(ctx context.Context, productID string) (*entity.Product, error)
	GetProductsForBuyer(ctx context.Context, limit int, cursor string) ([]entity.Product, error)
	CreateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, productID string) error
	UpdateProductForSeller(ctx context.Context, product *entity.Product, productID string, sellerID int64) error
	GetProductsForSeller(ctx context.Context, sellerId int64) ([]entity.Product, error)
	GetProductForSeller(ctx context.Context, productID string, sellerId int64) (*entity.Product, error)
}
