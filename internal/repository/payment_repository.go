package repository

import "github.com/febry3/gamingin/internal/entity"

type PaymentRepository interface {
	Create(payment *entity.Payment) error
	FindByOrderID(orderID string) (*entity.Payment, error)
	FindByGatewayTransactionID(txID string) (*entity.Payment, error)
	FindExpiredPending() ([]entity.Payment, error)
	Update(payment *entity.Payment) error
	UpdateStatus(paymentID, status string) error
}
