package entity

import "time"

type Order struct {
	ID                string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderNumber       string    `json:"order_number" gorm:"not null;uniqueIndex"`
	UserID            int64     `json:"user_id" gorm:"not null"`
	GroupBuySessionID *string   `json:"group_buy_session_id" gorm:"type:uuid;default:null"`
	SellerID          int64     `json:"seller_id" gorm:"not null"`
	Subtotal          float64   `json:"subtotal" gorm:"not null"`
	DeliveryCharge    float64   `json:"delivery_charge" gorm:"not null"`
	TotalAmount       float64   `json:"total_amount" gorm:"not null"`
	Status            string    `json:"status" gorm:"default:pending;check:status IN ('pending','processing','shipped','delivered','cancelled')"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (o *Order) TableName() string {
	return "orders"
}
