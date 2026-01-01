package repository

import "github.com/febry3/gamingin/internal/entity"

type OrderShippingDetailRepository interface {
	Create(detail *entity.OrderShippingDetail) error
	FindByOrderID(orderID string) (*entity.OrderShippingDetail, error)
}
