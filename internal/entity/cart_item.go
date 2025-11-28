package entity

import "time"

type CartItem struct {
	ID               string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CartID           string    `json:"cart_id" gorm:"type:uuid;not null"`
	ProductVariantID string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	Quantity         int       `json:"quantity" gorm:"default:1"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (ci *CartItem) TableName() string {
	return "cart_items"
}
