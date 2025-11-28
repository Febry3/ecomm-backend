package entity

import "time"

type OrderAdjustment struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID     string    `json:"order_id" gorm:"type:uuid;not null"`
	Description string    `json:"description" gorm:"not null"`
	Amount      float64   `json:"amount" gorm:"not null"`
	SourceType  string    `json:"source_type" gorm:"not null;check:source_type IN ('coupon','group_buy')"`
	SourceID    string    `json:"source_id" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (oa *OrderAdjustment) TableName() string {
	return "order_adjustments"
}
