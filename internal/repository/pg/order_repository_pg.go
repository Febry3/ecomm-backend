package pg

import (
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type OrderRepositoryPg struct {
	db *gorm.DB
}

func NewOrderRepositoryPg(db *gorm.DB) repository.OrderRepository {
	return &OrderRepositoryPg{db: db}
}

func (r *OrderRepositoryPg) Create(order *entity.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepositoryPg) FindByID(orderID string) (*entity.Order, error) {
	var order entity.Order
	err := r.db.
		Preload("ProductVariant").
		Preload("ProductVariant.Product").
		Preload("ProductVariant.Product.ProductImages", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("Payment").
		Preload("ShippingDetail").
		Preload("Seller").
		First(&order, "id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepositoryPg) FindByOrderNumber(orderNumber string) (*entity.Order, error) {
	var order entity.Order
	err := r.db.
		Preload("ProductVariant").
		Preload("Payment").
		Preload("ShippingDetail").
		First(&order, "order_number = ?", orderNumber).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepositoryPg) FindByUserID(userID int64, limit, offset int) ([]entity.Order, int64, error) {
	var orders []entity.Order
	var total int64

	// Get total count
	if err := r.db.Model(&entity.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated orders
	err := r.db.
		Preload("ProductVariant").
		Preload("ProductVariant.Product").
		Preload("ProductVariant.Product.ProductImages", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("Payment").
		Preload("Seller").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepositoryPg) Update(order *entity.Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepositoryPg) UpdateStatus(orderID, status string) error {
	return r.db.Model(&entity.Order{}).Where("id = ?", orderID).Update("status", status).Error
}
