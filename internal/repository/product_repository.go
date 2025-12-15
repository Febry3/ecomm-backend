package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type ProductRepository interface {
	GetProductForBuyer(ctx context.Context, productID string) (*entity.Product, error)
	GetProductsForBuyer(ctx context.Context) ([]entity.Product, error)
	CreateProduct(ctx context.Context, product *entity.Product) error
	UpdateProduct(ctx context.Context, product *entity.Product, productID string) error
	DeleteProduct(ctx context.Context, productID string) error
	GetProductsForSeller(ctx context.Context, sellerId int64) ([]entity.Product, error)
	GetProductForSeller(ctx context.Context, productID string, sellerId int64) (*entity.Product, error)
}
