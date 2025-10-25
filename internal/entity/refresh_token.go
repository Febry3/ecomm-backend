package entity

import "time"

type RefreshToken struct {
	TokenId          string    `json:"token_id" gorm:"primaryKey;"`
	UserId           int64     `json:"user_id" gorm:"foreignKey:UserId;references:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RefreshTokenHash string    `json:"refresh_token" gorm:"not null;size:255;uniqueIndex"`
	Role             string    `json:"role" gorm:"not null;default:user;check:role IN ('user','seller','admin')"`
	DeviceInfo       string    `json:"device_info" gorm:"size:255"`
	IpAddress        string    `json:"ip_address" gorm:"size:45"`
	Revoked          bool      `json:"revoked" gorm:"not null;default:false;index"`
	ExpiredAt        time.Time `json:"expired_at" gorm:"not null;index"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (r *RefreshToken) TableName() string {
	return "refresh_tokens"
}
