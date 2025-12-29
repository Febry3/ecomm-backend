package dto

import (
	"time"

	"github.com/febry3/gamingin/internal/entity"
)

type GroupBuySessionRequest struct {
	ProductVariantID string                `json:"product_variant_id" binding:"required"`
	MinParticipants  int                   `json:"min_participants" binding:"required,min=1"`
	MaxParticipants  int                   `json:"max_participants" binding:"required,min=1"`
	MaxQuantity      int                   `json:"max_quantity" binding:"required,min=1"`
	ExpiresAt        time.Time             `json:"expires_at" binding:"required"`
	Tiers            []GroupBuyTierRequest `json:"tiers" binding:"required,min=1,dive"`
}

type GroupBuyTierRequest struct {
	ParticipantThreshold int `json:"participant_threshold" binding:"required"`
	DiscountPercentage   int `json:"discount_percentage" binding:"required"`
}

type GroupBuySessionResponse struct {
	ID               string                 `json:"id"`
	SessionCode      string                 `json:"session_code"`
	ProductVariantID string                 `json:"product_variant_id"`
	SellerID         int64                  `json:"seller_id"`
	MinParticipants  int                    `json:"min_participants"`
	MaxParticipants  int                    `json:"max_participants"`
	Status           string                 `json:"status"`
	MaxQuantity      int64                  `json:"max_quantity"`
	ExpiresAt        time.Time              `json:"expires_at"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	ProductVariant   ProductVariantResponse `json:"product_variant,omitempty"`
	Tiers            []entity.GroupBuyTier  `json:"tiers"`
}

type ChangeStatusRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

type CreateBuyerGroupSessionRequest struct {
	OrganizerUserID  int64  `json:"organizer_user_id" binding:"required"`
	ProductVariantID string `json:"product_variant_id" binding:"required"`
	Title            string `json:"title" binding:"required"`
}

type GetBuyerGroupSessionResponse struct {
	Session        *entity.BuyerGroupSession `json:"buyer_group_session"`
	Address        []entity.Address          `json:"address"`
	ProductVariant *entity.ProductVariant    `json:"product_variant"`
	ProductSession *entity.GroupBuySession   `json:"product_session"`
}
