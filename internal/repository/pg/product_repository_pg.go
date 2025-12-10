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

func (p *ProductRepositoryPg) CreateProduct(ctx context.Context, product *entity.Product) error {
	db := TxFromContext(ctx, p.db)
	return db.Create(product).Error
}

func (p *ProductRepositoryPg) DeleteProduct(ctx context.Context, productID string) error {
	return p.db.WithContext(ctx).Delete(&entity.Product{}, productID).Error
}

func (p *ProductRepositoryPg) GetProduct(ctx context.Context, productID string) (*entity.Product, error) {
	var product entity.Product
	err := p.db.WithContext(ctx).First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

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
