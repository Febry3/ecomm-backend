package entity

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID            string           `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SellerID      int64            `json:"seller_id" gorm:"not null"`
	Title         string           `json:"title" gorm:"not null"`
	Slug          string           `json:"slug" gorm:"not null;uniqueIndex"`
	Description   datatypes.JSON   `json:"description" gorm:"type:jsonb"`
	CategoryID    int64            `json:"category_id" gorm:"default:null"`
	Badge         string           `json:"badge"`
	IsActive      bool             `json:"is_active" gorm:"default:true"`
	Status        string           `json:"status" gorm:"default:pending;check:status IN ('pending','approved','rejected')"`
	CreatedAt     time.Time        `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
	Variants      []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID;references:ID"`
	ProductImages []ProductImage   `json:"product_images,omitempty" gorm:"foreignKey:ProductID;references:ID"`
}

func (p *Product) TableName() string {
	return "products"
}
