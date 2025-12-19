package entity

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID            string           `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SellerID      int64            `json:"seller_id,omitempty" gorm:"not null"`
	Title         string           `json:"title,omitempty" gorm:"not null"`
	Slug          string           `json:"slug,omitempty" gorm:"not null;uniqueIndex"`
	Description   datatypes.JSON   `json:"description,omitempty" gorm:"type:jsonb"`
	CategoryID    int64            `json:"category_id,omitempty" gorm:"default:null"`
	Badge         string           `json:"badge,omitempty"`
	IsActive      bool             `json:"is_active,omitempty" gorm:"default:true"`
	Status        string           `json:"status,omitempty" gorm:"default:pending;check:status IN ('pending','approved','rejected')"`
	CreatedAt     time.Time        `json:"-" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt     time.Time        `json:"-" gorm:"autoUpdateTime;type:timestamptz"`
	Variants      []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID;references:ID"`
	ProductImages []ProductImage   `json:"product_images,omitempty" gorm:"foreignKey:ProductID;references:ID"`
}

func (p *Product) TableName() string {
	return "products"
}
