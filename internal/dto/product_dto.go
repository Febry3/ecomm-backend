package dto

import "encoding/json"

// CreateProductRequest is the main request for creating a product with variants
type CreateProductRequest struct {
	Title       string                  `json:"title" validate:"required"`
	Slug        string                  `json:"slug" validate:"required"`
	Description json.RawMessage         `json:"description"`
	CategoryID  string                  `json:"category_id,omitempty"`
	Badge       string                  `json:"badge,omitempty"`
	IsActive    bool                    `json:"is_active"`
	Variants    []ProductVariantRequest `json:"variants" validate:"required,min=1,dive"`
}

// ProductVariantRequest represents a product variant in the create request
type ProductVariantRequest struct {
	Sku      string                      `json:"sku" validate:"required"`
	Name     string                      `json:"name" validate:"required"`
	Price    float64                     `json:"price" validate:"required,gt=0"`
	IsActive bool                        `json:"is_active"`
	Stock    *ProductVariantStockRequest `json:"stock,omitempty"`
}

// ProductVariantStockRequest represents stock info for a variant
type ProductVariantStockRequest struct {
	CurrentStock      int `json:"current_stock"`
	ReservedStock     int `json:"reserved_stock"`
	LowStockThreshold int `json:"low_stock_threshold"`
}

// UpdateProductRequest for updating product details
type UpdateProductRequest struct {
	Title       string          `json:"title"`
	Slug        string          `json:"slug"`
	Description json.RawMessage `json:"description"`
	CategoryID  string          `json:"category_id,omitempty"`
	Badge       string          `json:"badge,omitempty"`
	IsActive    *bool           `json:"is_active,omitempty"`
}

// ProductResponse represents the response when returning a product
type ProductResponse struct {
	ID          string                   `json:"id"`
	SellerID    int64                    `json:"seller_id"`
	Title       string                   `json:"title"`
	Slug        string                   `json:"slug"`
	Description json.RawMessage          `json:"description"`
	CategoryID  string                   `json:"category_id,omitempty"`
	Badge       string                   `json:"badge,omitempty"`
	IsActive    bool                     `json:"is_active"`
	Status      string                   `json:"status"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
}

// ProductVariantResponse represents a variant in the response
type ProductVariantResponse struct {
	ID        string                       `json:"id"`
	ProductID string                       `json:"product_id"`
	Sku       string                       `json:"sku"`
	Name      string                       `json:"name"`
	Price     float64                      `json:"price"`
	IsActive  bool                         `json:"is_active"`
	CreatedAt string                       `json:"created_at"`
	UpdatedAt string                       `json:"updated_at"`
	Stock     *ProductVariantStockResponse `json:"stock,omitempty"`
}

// ProductVariantStockResponse represents stock info in the response
type ProductVariantStockResponse struct {
	CurrentStock      int    `json:"current_stock"`
	ReservedStock     int    `json:"reserved_stock"`
	LowStockThreshold int    `json:"low_stock_threshold"`
	LastUpdated       string `json:"last_updated"`
}
