package entity

import "time"

type GroupBuySession struct {
	ID                 string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionCode        string    `json:"session_code" gorm:"not null;uniqueIndex"`
	ProductVariantID   string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	OrganizerID        int64     `json:"organizer_id" gorm:"not null"`
	MinParticipants    int       `json:"min_participants" gorm:"not null"`
	MaxParticipants    int       `json:"max_participants" gorm:"not null"`
	DiscountPercentage float64   `json:"discount_percentage" gorm:"not null"`
	Status             string    `json:"status" gorm:"default:active;check:status IN ('active','completed','cancelled')"`
	ExpiresAt          time.Time `json:"expires_at" gorm:"not null;type:timestamptz"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (gbs *GroupBuySession) TableName() string {
	return "group_buy_sessions"
}
