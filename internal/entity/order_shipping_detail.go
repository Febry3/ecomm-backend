package entity

type OrderShippingDetail struct {
	ID           string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID      string `json:"order_id" gorm:"type:uuid;not null"`
	FullName     string `json:"full_name" gorm:"not null"`
	Phone        string `json:"phone" gorm:"not null"`
	AddressLine1 string `json:"address_line1" gorm:"not null"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city" gorm:"not null"`
	State        string `json:"state" gorm:"not null"`
	PostalCode   string `json:"postal_code" gorm:"not null"`
	Country      string `json:"country" gorm:"not null"`
}

func (osd *OrderShippingDetail) TableName() string {
	return "order_shipping_details"
}
