package entity

import (
	"time"
)

type RefreshToken struct {
	TokenId    string    `json:"token_id" gorm:"primaryKey;"`
	UserId     int64     `json:"user_id" gorm:"not null;uniqueIndex:ux_refresh_tokens_user_device"`
	TokenHash  string    `json:"token_hash" gorm:"not null;size:255;uniqueIndex"`
	Role       string    `json:"role" gorm:"not null;default:user;check:role IN ('user','seller','admin')"`
	IsRevoked  bool      `json:"is_revoked" gorm:"not null;default:false;index"`
	DeviceInfo string    `json:"device_info" gorm:"size:255;uniqueIndex:ux_refresh_tokens_user_device"`
	ExpiresAt  time.Time `json:"expired_at" gorm:"not null;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (r *RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (r *RefreshToken) IsExpired() bool {
	return r.ExpiresAt.Before(time.Now())
}
