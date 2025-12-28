package entity

import "time"

type BuyerGroupMember struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionID string    `json:"session_id" gorm:"type:uuid;not null;index"`
	UserID    int64     `json:"user_id" gorm:"not null;index"`
	OrderID   *string   `json:"order_id" gorm:"type:uuid;default:null"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Status    string    `json:"status" gorm:"default:joined;check:status IN ('joined','paid','cancelled')"`
	JoinedAt  time.Time `json:"joined_at" gorm:"autoCreateTime;type:timestamptz"`

	// Relationships
	Session *BuyerGroupSession `json:"session,omitempty" gorm:"foreignKey:SessionID;references:ID"`
	User    *User              `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Order   *Order             `json:"order,omitempty" gorm:"foreignKey:OrderID;references:ID"`
}

func (bgm *BuyerGroupMember) TableName() string {
	return "buyer_group_members"
}

// IsPaid checks if the member has completed payment
func (bgm *BuyerGroupMember) IsPaid() bool {
	return bgm.Status == "paid"
}
