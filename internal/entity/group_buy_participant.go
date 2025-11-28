package entity

import "time"

type GroupBuyParticipant struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SessionID string    `json:"session_id" gorm:"type:uuid;not null"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	JoinedAt  time.Time `json:"joined_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (gbp *GroupBuyParticipant) TableName() string {
	return "group_buy_participants"
}
