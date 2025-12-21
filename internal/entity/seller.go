package entity

import "time"

type Seller struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int64     `json:"user_id,omitempty" gorm:"not null;uniqueIndex"`
	StoreName     string    `json:"store_name" gorm:"not null"`
	StoreSlug     string    `json:"store_slug" gorm:"not null;uniqueIndex"`
	Description   string    `json:"description,omitempty" gorm:"type:text"`
	LogoURL       string    `json:"logo_url"`
	BusinessEmail string    `json:"business_email,omitempty"`
	BusinessPhone string    `json:"business_phone,omitempty"`
	Status        string    `json:"status,omitempty" gorm:"default:pending;check:status IN ('pending','approved','suspended')"`
	IsVerified    bool      `json:"is_verified,omitempty" gorm:"default:false"`
	AverageRating float64   `json:"average_rating,omitempty" gorm:"default:0"`
	TotalSales    int       `json:"total_sales,omitempty" gorm:"default:0"`
	CreatedAt     time.Time `json:"-" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt     time.Time `json:"-" gorm:"autoUpdateTime;type:timestamptz"`
}

func (s *Seller) TableName() string {
	return "sellers"
}
