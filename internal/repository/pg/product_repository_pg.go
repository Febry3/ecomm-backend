package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type ProductRepositoryPg struct {
	db *gorm.DB
}

func NewProductRepositoryPg(db *gorm.DB) repository.ProductRepository {
	return &ProductRepositoryPg{db: db}
}

// CreateProduct implements repository.ProductRepository.
func (p *ProductRepositoryPg) CreateProduct(ctx context.Context, product *entity.Product) error {
	return p.db.WithContext(ctx).Create(product).Error
}

// DeleteProduct implements repository.ProductRepository.
func (p *ProductRepositoryPg) DeleteProduct(ctx context.Context, productID string) error {
	return p.db.WithContext(ctx).Delete(&entity.Product{}, productID).Error
}

// GetProduct implements repository.ProductRepository.
func (p *ProductRepositoryPg) GetProduct(ctx context.Context, productID string) (*entity.Product, error) {
	var product entity.Product
	err := p.db.WithContext(ctx).First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProducts implements repository.ProductRepository.
func (p *ProductRepositoryPg) GetProducts(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product
	err := p.db.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct implements repository.ProductRepository.
func (p *ProductRepositoryPg) UpdateProduct(ctx context.Context, product *entity.Product, productID string) error {
	return p.db.WithContext(ctx).Save(product).Error
}
