package entity

import "time"

type Payment struct {
	ID                   string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID              string    `json:"order_id" gorm:"type:uuid;not null"`
	Amount               float64   `json:"amount" gorm:"not null"`
	Status               string    `json:"status" gorm:"default:pending;check:status IN ('pending','succeeded','failed','refunded')"`
	PaymentMethod        string    `json:"payment_method" gorm:"not null"`
	GatewayTransactionID string    `json:"gateway_transaction_id"`
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (p *Payment) TableName() string {
	return "payments"
}
