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
	return p.db.WithContext(ctx).Where("id = ?", productID).Delete(&entity.Product{}).Error
}

func (p *ProductRepositoryPg) GetProductForBuyer(ctx context.Context, productID string) (*entity.Product, error) {
	var product entity.Product
	err := p.db.WithContext(ctx).
		Preload("Variants").
		Preload("Variants.Stock").
		Preload("ProductImages").
		Preload("Seller", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "store_name", "store_slug", "logo_url")
		}).
		Where("id = ?", productID).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *ProductRepositoryPg) GetProductsForBuyer(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product
	err := p.db.WithContext(ctx).Preload("Variants").
		Preload("Variants.Stock").
		Preload("ProductImages").
		Preload("Seller", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "store_name", "store_slug", "logo_url")
		}).Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductRepositoryPg) UpdateProductForSeller(ctx context.Context, product *entity.Product, productID string, sellerID int64) error {
	db := TxFromContext(ctx, p.db)
	return db.Model(&entity.Product{}).Where("id = ? and seller_id = ?", productID, sellerID).Updates(product).Error
}

func (p *ProductRepositoryPg) GetProductForSeller(ctx context.Context, productID string, sellerId int64) (*entity.Product, error) {
	var products entity.Product
	err := p.db.WithContext(ctx).Where("seller_id = ? and id = ?", sellerId, productID).Preload("Variants").Preload("Variants.Stock").Preload("ProductImages").Order("created_at DESC").First(&products).Error
	if err != nil {
		return nil, err
	}
	return &products, nil
}

func (p *ProductRepositoryPg) GetProductsForSeller(ctx context.Context, sellerId int64) ([]entity.Product, error) {
	var products []entity.Product
	err := p.db.WithContext(ctx).Where("seller_id = ?", sellerId).Preload("ProductImages").Preload("Variants").Preload("Variants.Stock").Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
