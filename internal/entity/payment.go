package entity

import "time"

type Payment struct {
	ID                   string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID              string     `json:"order_id" gorm:"type:uuid;not null"`
	Amount               float64    `json:"amount" gorm:"not null"`
	Status               string     `json:"status" gorm:"default:pending"`
	PaymentMethod        string     `json:"payment_method" gorm:"default:bank_transfer"`
	BankCode             string     `json:"bank_code" gorm:"not null"`
	VANumber             string     `json:"va_number,omitempty"`
	BillKey              string     `json:"bill_key,omitempty"`    // For Mandiri Bill Payment
	BillerCode           string     `json:"biller_code,omitempty"` // For Mandiri Bill Payment
	GatewayTransactionID string     `json:"gateway_transaction_id,omitempty"`
	ExpiredAt            time.Time  `json:"expired_at" gorm:"not null;type:timestamptz"`
	PaidAt               *time.Time `json:"paid_at,omitempty" gorm:"type:timestamptz"`
	CreatedAt            time.Time  `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`

	// Relationships
	Order *Order `json:"order,omitempty" gorm:"foreignKey:OrderID;references:ID"`
}

func (p *Payment) TableName() string {
	return "payments"
}

// Payment status constants (matching Midtrans statuses)
const (
	PaymentStatusPending    = "pending"
	PaymentStatusSettlement = "settlement"
	PaymentStatusExpire     = "expire"
	PaymentStatusCancel     = "cancel"
	PaymentStatusDeny       = "deny"
)
