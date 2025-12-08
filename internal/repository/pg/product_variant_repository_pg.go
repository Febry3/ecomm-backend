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

// CreateProductVariant implements repository.ProductVariantRepository.
func (p *ProductVariantRepository) CreateProductVariant(ctx context.Context, productVariant *entity.ProductVariant) error {
	return p.db.WithContext(ctx).Create(productVariant).Error
}

// DeleteProductVariant implements repository.ProductVariantRepository.
func (p *ProductVariantRepository) DeleteProductVariant(ctx context.Context, productVariantID string) error {
	return p.db.WithContext(ctx).Delete(&entity.ProductVariant{}, productVariantID).Error
}

// GetProductVariant implements repository.ProductVariantRepository.
func (p *ProductVariantRepository) GetProductVariant(ctx context.Context, productVariantID string) (*entity.ProductVariant, error) {
	var productVariant entity.ProductVariant
	err := p.db.WithContext(ctx).First(&productVariant, productVariantID).Error
	if err != nil {
		return nil, err
	}
	return &productVariant, nil
}

// GetProductVariants implements repository.ProductVariantRepository.
func (p *ProductVariantRepository) GetProductVariants(ctx context.Context) ([]entity.ProductVariant, error) {
	var productVariants []entity.ProductVariant
	err := p.db.WithContext(ctx).Find(&productVariants).Error
	if err != nil {
		return nil, err
	}
	return productVariants, nil
}

// UpdateProductVariant implements repository.ProductVariantRepository.
func (p *ProductVariantRepository) UpdateProductVariant(ctx context.Context, productVariant *entity.ProductVariant, productVariantID string) error {
	return p.db.WithContext(ctx).Save(productVariant).Error
}
