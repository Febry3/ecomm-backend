package entity

import "time"

type User struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Username      string         `json:"username,omitempty" gorm:"default:null;uniqueIndex"`
	FirstName     string         `json:"first_name,omitempty" gorm:"default:null"`
	LastName      string         `json:"last_name,omitempty" gorm:"default:null"`
	PhoneNumber   string         `json:"phone_number,omitempty" gorm:"default:null;uniqueIndex"`
	Email         string         `json:"email,omitempty" gorm:"not null;uniqueIndex"`
	Role          string         `json:"role,omitempty" gorm:"type:text;check:role IN ('user','seller','admin');default:user;not null"`
	ProfileUrl    string         `json:"profile_url,omitempty" gorm:"default:null"`
	CreatedAt     *time.Time     `json:"created_at,omitempty" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty" gorm:"autoUpdateTime;type:timestamptz"`
	AuthProviders []AuthProvider `json:"auth_providers,omitempty" gorm:"foreignKey:UserId"`
	RefreshTokens []RefreshToken `json:"refresh_tokens,omitempty" gorm:"foreignKey:UserId"`
}

func (r *User) TableName() string {
	return "users"
}
