package entity

import "time"

type UserFavorite struct {
	ID               string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID           int64     `json:"user_id" gorm:"not null"`
	ProductVariantID string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (uf *UserFavorite) TableName() string {
	return "user_favorites"
}
