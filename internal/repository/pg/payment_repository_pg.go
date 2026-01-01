package pg

import (
	"time"

	"github.com/febry3/gamingin/internal/entity"
	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type PaymentRepositoryPg struct {
	db *gorm.DB
}

func NewPaymentRepositoryPg(db *gorm.DB) repository.PaymentRepository {
	return &PaymentRepositoryPg{db: db}
}

func (r *PaymentRepositoryPg) Create(payment *entity.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepositoryPg) FindByOrderID(orderID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.First(&payment, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryPg) FindByGatewayTransactionID(txID string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.Preload("Order").First(&payment, "gateway_transaction_id = ?", txID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryPg) FindExpiredPending() ([]entity.Payment, error) {
	var payments []entity.Payment
	err := r.db.
		Preload("Order").
		Where("status = ? AND expired_at < ?", entity.PaymentStatusPending, time.Now()).
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepositoryPg) Update(payment *entity.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepositoryPg) UpdateStatus(paymentID, status string) error {
	return r.db.Model(&entity.Payment{}).Where("id = ?", paymentID).Update("status", status).Error
}
