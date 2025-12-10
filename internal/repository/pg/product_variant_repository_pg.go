package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type ProductVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepositoryPg(db *gorm.DB) repository.ProductVariantRepository {
	return &ProductVariantRepository{db: db}
}

func (p *ProductVariantRepository) CreateProductVariant(ctx context.Context, productVariant *entity.ProductVariant) error {
	db := TxFromContext(ctx, p.db)
	return db.Create(productVariant).Error
}

func (p *ProductVariantRepository) DeleteProductVariant(ctx context.Context, productVariantID string) error {
	return p.db.WithContext(ctx).Delete(&entity.ProductVariant{}, productVariantID).Error
}

func (p *ProductVariantRepository) GetProductVariant(ctx context.Context, productVariantID string) (*entity.ProductVariant, error) {
	var productVariant entity.ProductVariant
	err := p.db.WithContext(ctx).First(&productVariant, productVariantID).Error
	if err != nil {
		return nil, err
	}
	return &productVariant, nil
}

func (p *ProductVariantRepository) GetProductVariants(ctx context.Context) ([]entity.ProductVariant, error) {
	var productVariants []entity.ProductVariant
	err := p.db.WithContext(ctx).Find(&productVariants).Error
	if err != nil {
		return nil, err
	}
	return productVariants, nil
}

func (p *ProductVariantRepository) UpdateProductVariant(ctx context.Context, productVariant *entity.ProductVariant, productVariantID string) error {
	return p.db.WithContext(ctx).Save(productVariant).Error
}
