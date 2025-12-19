package entity

import "time"

type ProductImage struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID    string    `json:"product_id" gorm:"type:uuid;not null"`
	ImageURL     string    `json:"image_url" gorm:"not null"`
	AltText      string    `json:"alt_text,omitempty"`
	DisplayOrder int       `json:"display_order,omitempty" gorm:"default:0"`
	IsPrimary    bool      `json:"is_primary,omitempty" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"autoCreateTime;type:timestamptz"`
}

func (pi *ProductImage) TableName() string {
	return "product_images"
}
