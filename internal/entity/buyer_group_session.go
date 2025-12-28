package entity

import "time"

type BuyerGroupSession struct {
	ID                  string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionCode         string    `json:"session_code" gorm:"not null;uniqueIndex"`
	OrganizerUserID     int64     `json:"organizer_user_id" gorm:"not null"`
	ProductVariantID    string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	Title               string    `json:"title" gorm:"default:null"`
	CurrentParticipants int       `json:"current_participants" gorm:"default:1"`
	Status              string    `json:"status" gorm:"default:open;check:status IN ('open','locked','completed','cancelled','expired')"`
	ExpiresAt           time.Time `json:"expires_at" gorm:"not null;type:timestamptz"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`

	// Relationships
	Organizer      *User              `json:"organizer,omitempty" gorm:"foreignKey:OrganizerUserID;references:ID"`
	ProductVariant *ProductVariant    `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID;references:ID"`
	Members        []BuyerGroupMember `json:"members,omitempty" gorm:"foreignKey:SessionID;references:ID"`
}

func (bgs *BuyerGroupSession) TableName() string {
	return "buyer_group_sessions"
}

func (bgs *BuyerGroupSession) IsExpired() bool {
	return bgs.ExpiresAt.Before(time.Now())
}
