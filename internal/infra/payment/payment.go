package payment

import (
	"context"
	"time"
)

// VAPaymentResult represents the result of creating a VA payment
type VAPaymentResult struct {
	TransactionID string
	OrderID       string
	Bank          string
	VANumber      string
	BillKey       string // For Mandiri Bill Payment
	BillerCode    string // For Mandiri Bill Payment
	GrossAmount   float64
	Status        string
	ExpiredAt     time.Time
}

// PaymentStatusResult represents the result of checking payment status
type PaymentStatusResult struct {
	TransactionID string
	OrderID       string
	Status        string // pending, settlement, expire, cancel, deny
	PaymentType   string
	GrossAmount   float64
	PaidAt        *time.Time
}

// PaymentGateway defines the interface for payment gateway operations
type PaymentGateway interface {
	// ChargeVA creates a Virtual Account payment
	ChargeVA(ctx context.Context, orderID string, amount int64, bankCode string) (*VAPaymentResult, error)

	// GetTransactionStatus checks the current status of a transaction
	GetTransactionStatus(ctx context.Context, orderID string) (*PaymentStatusResult, error)

	// VerifySignature validates the webhook notification signature
	VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool

	// CancelTransaction cancels a pending transaction
	CancelTransaction(ctx context.Context, orderID string) error
}
