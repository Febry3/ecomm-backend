package repository

import "github.com/febry3/gamingin/internal/entity"

type OrderRepository interface {
	Create(order *entity.Order) error
	FindByID(orderID string) (*entity.Order, error)
	FindByOrderNumber(orderNumber string) (*entity.Order, error)
	FindByUserID(userID int64, limit, offset int) ([]entity.Order, int64, error)
	Update(order *entity.Order) error
	UpdateStatus(orderID, status string) error
}
