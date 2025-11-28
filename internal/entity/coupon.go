package entity

import "time"

type Coupon struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Code              string    `json:"code" gorm:"not null;uniqueIndex"`
	DiscountType      string    `json:"discount_type" gorm:"not null;check:discount_type IN ('percentage','fixed')"`
	DiscountValue     float64   `json:"discount_value" gorm:"not null"`
	MinPurchaseAmount float64   `json:"min_purchase_amount" gorm:"default:0"`
	MaxDiscountAmount float64   `json:"max_discount_amount" gorm:"default:0"`
	ValidFrom         time.Time `json:"valid_from" gorm:"not null;type:timestamptz"`
	ValidUntil        time.Time `json:"valid_until" gorm:"not null;type:timestamptz"`
	UsageLimit        int       `json:"usage_limit" gorm:"default:0"`
	UsageCount        int       `json:"usage_count" gorm:"default:0"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (c *Coupon) TableName() string {
	return "coupons"
}
