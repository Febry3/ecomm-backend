package entity

import "time"

type User struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Username      string         `json:"username" gorm:"not null;uniqueIndex"`
	FirstName     string         `json:"first_name" gorm:"not null"`
	LastName      string         `json:"last_name" gorm:"not null"`
	PhoneNumber   string         `json:"phone_number" gorm:"not null;uniqueIndex"`
	Email         string         `json:"email" gorm:"not null;uniqueIndex"`
	Role          string         `json:"role" gorm:"type:text;check:role IN ('user','seller','admin');default:user;not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	AuthProviders []AuthProvider `gorm:"foreignKey:UserId"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserId"`
}

func (r *User) TableName() string {
	return "users"
}
