package entity

import "time"

type Order struct {
	ID                  string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderNumber         string    `json:"order_number" gorm:"not null;uniqueIndex"`
	UserID              int64     `json:"user_id" gorm:"not null"`
	BuyerGroupSessionID *string   `json:"buyer_group_session_id,omitempty" gorm:"type:uuid"`
	SellerID            int64     `json:"seller_id" gorm:"not null"`
	ProductVariantID    string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	Quantity            int       `json:"quantity" gorm:"not null;default:1"`
	PriceAtOrder        float64   `json:"price_at_order" gorm:"not null"`
	Subtotal            float64   `json:"subtotal" gorm:"not null"`
	DeliveryCharge      float64   `json:"delivery_charge" gorm:"default:0"`
	TotalAmount         float64   `json:"total_amount" gorm:"not null"`
	Status              string    `json:"status" gorm:"default:pending_payment"`
	AddressID           string    `json:"address_id" gorm:"type:uuid;not null"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`

	// Relationships
	User              *User                `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Seller            *Seller              `json:"seller,omitempty" gorm:"foreignKey:SellerID;references:ID"`
	ProductVariant    *ProductVariant      `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID;references:ID"`
	BuyerGroupSession *BuyerGroupSession   `json:"buyer_group_session,omitempty" gorm:"foreignKey:BuyerGroupSessionID;references:ID"`
	Payment           *Payment             `json:"payment,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	ShippingDetail    *OrderShippingDetail `json:"shipping_detail,omitempty" gorm:"foreignKey:OrderID;references:ID"`
}

func (o *Order) TableName() string {
	return "orders"
}

// Order status constants
const (
	OrderStatusPendingPayment = "pending_payment"
	OrderStatusPaid           = "paid"
	OrderStatusProcessing     = "processing"
	OrderStatusShipped        = "shipped"
	OrderStatusDelivered      = "delivered"
	OrderStatusCancelled      = "cancelled"
	OrderStatusExpired        = "expired"
)
