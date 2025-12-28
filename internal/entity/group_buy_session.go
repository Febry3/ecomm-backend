package entity

import (
	"time"
)

type GroupBuySession struct {
	ID                  string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionCode         string    `json:"session_code,omitempty" gorm:"not null;uniqueIndex"`
	ProductVariantID    string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	SellerID            int64     `json:"seller_id" gorm:"not null"`
	MinParticipants     int       `json:"min_participants,omitempty" gorm:"not null"`
	MaxParticipants     int       `json:"max_participants,omitempty" gorm:"not null"`
	CurrentParticipants int       `json:"current_participants,omitempty" gorm:"default:0"`
	FinalTierID         *string   `json:"final_tier_id,omitempty" gorm:"type:uuid;default:null"`
	Status              string    `json:"status,omitempty" gorm:"default:active;check:status IN ('active','completed','cancelled')"`
	MaxQuantity         int64     `json:"max_quantity,omitempty" gorm:"not null"`
	ExpiresAt           time.Time `json:"expires_at,omitempty" gorm:"not null;type:timestamptz"`
	CreatedAt           time.Time `json:"-" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt           time.Time `json:"-" gorm:"autoUpdateTime;type:timestamptz"`
	// ProductVariant is loaded from the other side, don't include here to avoid circular reference
	ProductVariant ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`
	// Seller           Seller         `json:"seller" gorm:"foreignKey:SellerID"`
	GroupBuyTiers []GroupBuyTier `json:"group_buy_tiers,omitempty" gorm:"foreignKey:GroupBuySessionID"`
}

func (gbs *GroupBuySession) TableName() string {
	return "group_buy_sessions"
}
