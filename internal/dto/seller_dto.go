package dto

type SellerRequest struct {
	StoreName     string `json:"store_name" form:"store_name" validate:"required"`
	StoreSlug     string `json:"store_slug" form:"store_slug" validate:"required"`
	Description   string `json:"description" form:"description" validate:"required"`
	BusinessEmail string `json:"business_email" form:"business_email" validate:"required,email"`
	BusinessPhone string `json:"business_phone" form:"business_phone" validate:"required"`
}

type UpdateSellerRequest struct {
	StoreName     string `json:"store_name" form:"store_name" validate:"required"`
	StoreSlug     string `json:"store_slug" form:"store_slug" validate:""`
	Description   string `json:"description" form:"description" validate:"required"`
	BusinessEmail string `json:"business_email" form:"business_email" validate:"required,email"`
	BusinessPhone string `json:"business_phone" form:"business_phone" validate:"required"`
}
