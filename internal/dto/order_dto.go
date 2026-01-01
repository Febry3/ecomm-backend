package dto

import "time"

// ========================================
// Request DTOs
// ========================================

// CreateOrderRequest for direct buy flow
type CreateOrderRequest struct {
	ProductVariantID string `json:"product_variant_id" validate:"required,uuid"`
	Quantity         int    `json:"quantity" validate:"required,min=1"`
	AddressID        string `json:"address_id" validate:"required,uuid"`
	BankCode         string `json:"bank_code" validate:"required,oneof=bca bni bri mandiri permata cimb"`
}

// CreateGroupBuyOrderRequest for group buy flow
type CreateGroupBuyOrderRequest struct {
	BuyerGroupSessionID string `json:"buyer_group_session_id" validate:"required,uuid"`
	AddressID           string `json:"address_id" validate:"required,uuid"`
	BankCode            string `json:"bank_code" validate:"required,oneof=bca bni bri mandiri permata cimb"`
}

// ========================================
// Response DTOs
// ========================================

// OrderResponse is the main order response
type OrderResponse struct {
	ID             string                  `json:"id"`
	OrderNumber    string                  `json:"order_number"`
	Status         string                  `json:"status"`
	Quantity       int                     `json:"quantity"`
	PriceAtOrder   float64                 `json:"price_at_order"`
	Subtotal       float64                 `json:"subtotal"`
	DeliveryCharge float64                 `json:"delivery_charge"`
	TotalAmount    float64                 `json:"total_amount"`
	Payment        *PaymentDetailResponse  `json:"payment,omitempty"`
	Product        *OrderProductResponse   `json:"product,omitempty"`
	ShippingDetail *ShippingDetailResponse `json:"shipping_detail,omitempty"`
	Seller         *OrderSellerResponse    `json:"seller,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
}

// PaymentDetailResponse contains VA payment info
type PaymentDetailResponse struct {
	ID         string     `json:"id"`
	BankCode   string     `json:"bank_code"`
	VANumber   string     `json:"va_number,omitempty"`
	BillKey    string     `json:"bill_key,omitempty"`    // For Mandiri
	BillerCode string     `json:"biller_code,omitempty"` // For Mandiri
	Amount     float64    `json:"amount"`
	Status     string     `json:"status"`
	ExpiredAt  time.Time  `json:"expired_at"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
}

// OrderProductResponse contains product info for order
type OrderProductResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	VariantID   string `json:"variant_id"`
	VariantName string `json:"variant_name"`
	ImageURL    string `json:"image_url,omitempty"`
}

// ShippingDetailResponse contains shipping address info
type ShippingDetailResponse struct {
	ReceiverName  string `json:"receiver_name"`
	Phone         string `json:"phone,omitempty"`
	StreetAddress string `json:"street_address"`
	Village       string `json:"village,omitempty"`
	District      string `json:"district,omitempty"`
	City          string `json:"city"`
	Province      string `json:"province"`
	PostalCode    string `json:"postal_code"`
}

// OrderSellerResponse contains seller info
type OrderSellerResponse struct {
	ID       int64  `json:"id"`
	ShopName string `json:"shop_name"`
}

// OrderListResponse for paginated order list
type OrderListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
}

// ========================================
// Midtrans Webhook DTOs
// ========================================

// MidtransNotification represents the webhook payload from Midtrans
type MidtransNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}
