package entity

import "time"

type InventoryLedger struct {
	ID               int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductVariantID string    `json:"product_variant_id" gorm:"type:uuid;not null"`
	QuantityChange   int       `json:"quantity_change" gorm:"not null"`
	Reason           string    `json:"reason" gorm:"not null"`
	OrderID          *string   `json:"order_id" gorm:"type:uuid;default:null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
}

func (il *InventoryLedger) TableName() string {
	return "inventory_ledgers"
}
