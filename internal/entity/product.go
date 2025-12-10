package entity

import "time"

type Product struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SellerID    int64     `json:"seller_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	CategoryID  string    `json:"category_id" gorm:"type:uuid"` // commented the not null
	Badge       string    `json:"badge"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Status      string    `json:"status" gorm:"default:pending;check:status IN ('pending','approved','rejected')"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (p *Product) TableName() string {
	return "products"
}
