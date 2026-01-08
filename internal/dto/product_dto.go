package dto

import (
	"encoding/json"

	"github.com/febry3/gamingin/internal/entity"
)

// CreateProductRequest is the main request for creating a product with variants
type CreateProductRequest struct {
	Title       string                  `json:"title" validate:"required"`
	Slug        string                  `json:"slug" validate:"required"`
	Description json.RawMessage         `json:"description"`
	CategoryID  int64                   `json:"category_id,omitempty"`
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

// ProductResponse represents the response when returning a product
type ProductResponse struct {
	ID          string                   `json:"id"`
	SellerID    int64                    `json:"seller_id"`
	Title       string                   `json:"title"`
	Slug        string                   `json:"slug"`
	Description json.RawMessage          `json:"description"`
	CategoryID  int64                    `json:"category_id,omitempty"`
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

// UpdateProductRequest for updating product details
type UpdateProductRequest struct {
	Title           string                        `json:"title"`
	Slug            string                        `json:"slug"`
	Description     json.RawMessage               `json:"description"`
	CategoryID      int64                         `json:"category_id,omitempty"`
	Badge           string                        `json:"badge,omitempty"`
	IsActive        *bool                         `json:"is_active,omitempty"`
	ProductVariants []UpdateProductVariantRequest `json:"variants,omitempty"`
}

type UpdateProductVariantRequest struct {
	ID        string                      `json:"id" validate:"required"`
	ProductID string                      `json:"product_id" validate:"required"`
	Sku       string                      `json:"sku" validate:"required"`
	Name      string                      `json:"name" validate:"required"`
	Price     float64                     `json:"price" validate:"required,gt=0"`
	IsActive  bool                        `json:"is_active"`
	Stock     *ProductVariantStockRequest `json:"stock,omitempty"`
}

type GetProductsResponse struct {
	Products []entity.Product `json:"products"`
	Cursor   string           `json:"cursor,omitempty"`
	HasMore  bool             `json:"has_more"`
}

func ToGetProductResponse(products []entity.Product, limit int) GetProductsResponse {
	if len(products) == 0 {
		return GetProductsResponse{
			Products: []entity.Product{},
			Cursor:   "",
			HasMore:  false,
		}
	}

	hasMore := len(products) > limit
	resultProducts := products
	if hasMore {
		resultProducts = products[:limit]
	}

	var cursor string
	if hasMore && len(resultProducts) > 0 {
		cursor = resultProducts[len(resultProducts)-1].CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z")
	}

	return GetProductsResponse{
		Products: resultProducts,
		Cursor:   cursor,
		HasMore:  hasMore,
	}
}

// ToProductResponse converts a Product entity and its variants to ProductResponse DTO
func ToProductResponse(product *entity.Product, variants []entity.ProductVariant) *ProductResponse {
	variantResponses := make([]ProductVariantResponse, 0, len(variants))
	for _, v := range variants {
		variantResponses = append(variantResponses, ToProductVariantResponse(&v))
	}

	return &ProductResponse{
		ID:          product.ID,
		SellerID:    product.SellerID,
		Title:       product.Title,
		Slug:        product.Slug,
		Description: json.RawMessage(product.Description),
		CategoryID:  product.CategoryID,
		Badge:       product.Badge,
		IsActive:    product.IsActive,
		Status:      product.Status,
		Variants:    variantResponses,
	}
}

// ToProductVariantResponse converts a ProductVariant entity to ProductVariantResponse DTO
func ToProductVariantResponse(variant *entity.ProductVariant) ProductVariantResponse {
	response := ProductVariantResponse{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		Sku:       variant.Sku,
		Name:      variant.Name,
		Price:     variant.Price,
		IsActive:  variant.IsActive,
		CreatedAt: variant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: variant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if variant.Stock != nil {
		response.Stock = ToProductVariantStockResponse(variant.Stock)
	}

	return response
}

// ToProductVariantStockResponse converts a ProductVariantStock entity to ProductVariantStockResponse DTO
func ToProductVariantStockResponse(stock *entity.ProductVariantStock) *ProductVariantStockResponse {
	return &ProductVariantStockResponse{
		CurrentStock:      stock.CurrentStock,
		ReservedStock:     stock.ReservedStock,
		LowStockThreshold: stock.LowStockThreshold,
		LastUpdated:       stock.LastUpdated.Format("2006-01-02T15:04:05Z07:00"),
	}
}
