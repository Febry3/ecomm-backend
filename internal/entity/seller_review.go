package entity

import "time"

type SellerReview struct {
	ID                  string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SellerID            int64     `json:"seller_id" gorm:"not null"`
	UserID              int64     `json:"user_id" gorm:"not null"`
	OrderID             string    `json:"order_id" gorm:"type:uuid;not null"`
	Rating              int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	ReviewText          string    `json:"review_text" gorm:"type:text"`
	CommunicationRating int       `json:"communication_rating" gorm:"not null;check:communication_rating >= 1 AND communication_rating <= 5"`
	ShippingRating      int       `json:"shipping_rating" gorm:"not null;check:shipping_rating >= 1 AND shipping_rating <= 5"`
	ProductRating       int       `json:"product_rating" gorm:"not null;check:product_rating >= 1 AND product_rating <= 5"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (sr *SellerReview) TableName() string {
	return "seller_reviews"
}
