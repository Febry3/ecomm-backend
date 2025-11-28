package entity

import "time"

type ProductReview struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProductID  string    `json:"product_id" gorm:"type:uuid;not null"`
	UserID     int64     `json:"user_id" gorm:"not null"`
	Rating     int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	ReviewText string    `json:"review_text" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (pr *ProductReview) TableName() string {
	return "product_reviews"
}
