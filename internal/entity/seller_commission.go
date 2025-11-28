package entity

import "time"

type SellerCommission struct {
	ID               string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SellerID         int64      `json:"seller_id" gorm:"not null"`
	OrderID          string     `json:"order_id" gorm:"type:uuid;not null"`
	OrderItemID      string     `json:"order_item_id" gorm:"type:uuid;not null"`
	SaleAmount       float64    `json:"sale_amount" gorm:"not null"`
	CommissionRate   float64    `json:"commission_rate" gorm:"not null"`
	CommissionAmount float64    `json:"commission_amount" gorm:"not null"`
	SellerEarnings   float64    `json:"seller_earnings" gorm:"not null"`
	Status           string     `json:"status" gorm:"default:pending;check:status IN ('pending','paid','held')"`
	PayoutID         *string    `json:"payout_id" gorm:"type:uuid;default:null"`
	PaidAt           *time.Time `json:"paid_at" gorm:"default:null;type:timestamptz"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (sc *SellerCommission) TableName() string {
	return "seller_commissions"
}
