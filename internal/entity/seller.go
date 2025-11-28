package entity

import "time"

type Seller struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int64     `json:"user_id" gorm:"not null;uniqueIndex"`
	StoreName     string    `json:"store_name" gorm:"not null"`
	StoreSlug     string    `json:"store_slug" gorm:"not null;uniqueIndex"`
	Description   string    `json:"description" gorm:"type:text"`
	LogoURL       string    `json:"logo_url"`
	BusinessEmail string    `json:"business_email"`
	BusinessPhone string    `json:"business_phone"`
	Status        string    `json:"status" gorm:"default:pending;check:status IN ('pending','approved','suspended')"`
	IsVerified    bool      `json:"is_verified" gorm:"default:false"`
	AverageRating float64   `json:"average_rating" gorm:"default:0"`
	TotalSales    int       `json:"total_sales" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (s *Seller) TableName() string {
	return "sellers"
}
