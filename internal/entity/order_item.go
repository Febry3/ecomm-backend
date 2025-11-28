package entity

import "time"

type OrderItem struct {
	ID               string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID          string    `json:"order_id" gorm:"type:uuid;not null"`
	ProductVariantID string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	Quantity         int       `json:"quantity" gorm:"not null"`
	PriceAtPurchase  float64   `json:"price_at_purchase" gorm:"not null"`
	TotalPrice       float64   `json:"total_price" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (oi *OrderItem) TableName() string {
	return "order_items"
}
