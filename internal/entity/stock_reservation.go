package entity

import "time"

type StockReservation struct {
	ID               string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductVariantID string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	UserID           int64     `json:"user_id" gorm:"not null"`
	OrderID          *string   `json:"order_id" gorm:"type:uuid;default:null"`
	Quantity         int       `json:"quantity" gorm:"not null"`
	Status           string    `json:"status" gorm:"default:pending;check:status IN ('pending','completed','expired')"`
	ExpiresAt        time.Time `json:"expires_at" gorm:"not null;type:timestamptz;index"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (sr *StockReservation) TableName() string {
	return "stock_reservations"
}

func (sr *StockReservation) IsExpired() bool {
	return sr.ExpiresAt.Before(time.Now())
}
