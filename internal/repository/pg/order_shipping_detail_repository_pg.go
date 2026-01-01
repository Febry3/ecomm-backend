package pg

import (
	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type OrderShippingDetailRepositoryPg struct {
	db *gorm.DB
}

func NewOrderShippingDetailRepositoryPg(db *gorm.DB) repository.OrderShippingDetailRepository {
	return &OrderShippingDetailRepositoryPg{db: db}
}

func (r *OrderShippingDetailRepositoryPg) Create(detail *entity.OrderShippingDetail) error {
	return r.db.Create(detail).Error
}

func (r *OrderShippingDetailRepositoryPg) FindByOrderID(orderID string) (*entity.OrderShippingDetail, error) {
	var detail entity.OrderShippingDetail
	err := r.db.First(&detail, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}
