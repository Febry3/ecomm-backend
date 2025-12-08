package entity

import "time"

type ProductVariantStock struct {
	ProductVariantID  string          `json:"product_variant_id" gorm:"primaryKey;type:uuid"`
	CurrentStock      int             `json:"current_stock" gorm:"default:0"`
	ReservedStock     int             `json:"reserved_stock" gorm:"default:0"`
	LowStockThreshold int             `json:"low_stock_threshold" gorm:"default:5"`
	LastUpdated       time.Time       `json:"last_updated" gorm:"autoUpdateTime;type:timestamptz"`
	ProductVariant    *ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID;references:ID"`
}

func (pvs *ProductVariantStock) TableName() string {
	return "product_variant_stocks"
}
