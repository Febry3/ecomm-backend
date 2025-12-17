package entity

import (
	"time"
)

type GroupBuySession struct {
	ID               string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionCode      string         `json:"session_code" gorm:"not null;uniqueIndex"`
	ProductVariantID string         `json:"product_variant_id" gorm:"type:uuid;not null"`
	SellerID         int64          `json:"seller_id" gorm:"not null"`
	MinParticipants  int            `json:"min_participants" gorm:"not null"`
	MaxParticipants  int            `json:"max_participants" gorm:"not null"`
	Status           string         `json:"status" gorm:"default:active;check:status IN ('active','completed','cancelled')"`
	MaxQuantity      int64          `json:"max_quantity" gorm:"not null"`
	ExpiresAt        time.Time      `json:"expires_at" gorm:"not null;type:timestamptz"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
	ProductVariant   ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`
	Seller           Seller         `json:"seller" gorm:"foreignKey:SellerID"`
	GroupBuyTiers    []GroupBuyTier `json:"group_buy_tiers,omitempty" gorm:"foreignKey:GroupBuySessionID"`
}

func (gbs *GroupBuySession) TableName() string {
	return "group_buy_sessions"
}
