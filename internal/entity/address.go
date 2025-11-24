package entity

import (
	"time"
)

type Address struct {
	AddressID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"address_id"`
	UserID    int64  `gorm:"not null;index" json:"user_id"`
	// Physical Address
	StreetAddress string `gorm:"type:text;not null" json:"street_address"` // Jalan + No Rumah
	// RT/RW are optional (omitempty) because some cluster housing doesn't use them
	RT string `gorm:"type:varchar(5)" json:"rt,omitempty"`
	RW string `gorm:"type:varchar(5)" json:"rw,omitempty"`
	// Regional Hierarchy
	Village    string `gorm:"type:varchar(100);not null" json:"village"`    // Kelurahan
	District   string `gorm:"type:varchar(100);not null" json:"district"`   // Kecamatan
	City       string `gorm:"type:varchar(100);not null" json:"city"`       // Kota/Kabupaten
	Province   string `gorm:"type:varchar(100);not null" json:"province"`   // Provinsi
	PostalCode string `gorm:"type:varchar(10);not null" json:"postal_code"` // Kode Pos
	// Notes (Patokan) - Optional
	Notes string `gorm:"type:text" json:"notes,omitempty"`
	// Metadata
	IsDefault bool      `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
