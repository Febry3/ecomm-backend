package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type ProductImageRepositoryPg struct {
	db *gorm.DB
}

func NewProductImageRepositoryPg(db *gorm.DB) repository.ProductImageRepository {
	return &ProductImageRepositoryPg{db: db}
}

func (p *ProductImageRepositoryPg) CreateProductImage(ctx context.Context, productImage *entity.ProductImage) error {
	db := TxFromContext(ctx, p.db)
	return db.Create(productImage).Error
}
