package repository

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
)

type ProductRepository interface {
	GetProduct(ctx context.Context, productID string) (*entity.Product, error)
	GetProducts(ctx context.Context) ([]entity.Product, error)
	CreateProduct(ctx context.Context, product *entity.Product) error
	UpdateProduct(ctx context.Context, product *entity.Product, productID string) error
	DeleteProduct(ctx context.Context, productID string) error
}
