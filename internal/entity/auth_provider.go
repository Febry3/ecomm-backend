package entity

import "database/sql"

type AuthProvider struct {
	AuthProviderID int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId         int64          `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;uniqueIndex:idx_user_provider"`
	Provider       string         `json:"provider" gorm:"not null;size:50;uniqueIndex:idx_user_provider"`
	ProviderId     string         `json:"provider_id" gorm:"not null;size:255;uniqueIndex"`
	Password       sql.NullString `json:"-" gorm:"default:null"`
}

func (r *AuthProvider) TableName() string {
	return "auth_providers"
}
