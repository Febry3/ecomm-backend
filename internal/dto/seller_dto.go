package dto

type SellerRequest struct {
	StoreName     string `json:"store_name" validate:"required"`
	StoreSlug     string `json:"store_slug" validate:"required"`
	Description   string `json:"description" validate:"required"`
	LogoURL       string `json:"logo_url" validate:"required"`
	BusinessEmail string `json:"business_email" validate:"required,email"`
	BusinessPhone string `json:"business_phone" validate:"required,number,min=8,max=12"`
}
