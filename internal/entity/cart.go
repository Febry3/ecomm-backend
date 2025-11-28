package entity

import "time"

type Cart struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    int64     `json:"user_id" gorm:"not null"`
	SessionID string    `json:"session_id" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (c *Cart) TableName() string {
	return "carts"
}
