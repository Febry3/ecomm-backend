package entity

import "time"

type ProductVariant struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID string    `json:"product_id" gorm:"type:uuid;not null"`
	Sku       string    `json:"sku" gorm:"not null;uniqueIndex"`
	Name      string    `json:"name" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (pv *ProductVariant) TableName() string {
	return "product_variants"
}
