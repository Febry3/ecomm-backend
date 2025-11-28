package entity

import "time"

type SellerPayout struct {
	ID            string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SellerID      int64      `json:"seller_id" gorm:"not null"`
	Amount        float64    `json:"amount" gorm:"not null"`
	Status        string     `json:"status" gorm:"default:pending;check:status IN ('pending','processing','completed','failed')"`
	PaymentMethod string     `json:"payment_method" gorm:"not null;check:payment_method IN ('bank_transfer','paypal')"`
	TransactionID string     `json:"transaction_id"`
	ProcessedAt   *time.Time `json:"processed_at" gorm:"default:null;type:timestamptz"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (sp *SellerPayout) TableName() string {
	return "seller_payouts"
}
