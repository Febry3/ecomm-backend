package entity

type OrderShippingDetail struct {
	ID            string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID       string `json:"order_id" gorm:"type:uuid;not null;uniqueIndex"`
	ReceiverName  string `json:"receiver_name" gorm:"not null"`
	Phone         string `json:"phone,omitempty"`
	StreetAddress string `json:"street_address" gorm:"type:text;not null"`
	RT            string `json:"rt,omitempty" gorm:"type:varchar(5)"`
	RW            string `json:"rw,omitempty" gorm:"type:varchar(5)"`
	Village       string `json:"village,omitempty" gorm:"type:varchar(100)"`
	District      string `json:"district,omitempty" gorm:"type:varchar(100)"`
	City          string `json:"city" gorm:"type:varchar(100);not null"`
	Province      string `json:"province" gorm:"type:varchar(100);not null"`
	PostalCode    string `json:"postal_code" gorm:"type:varchar(10);not null"`
	Notes         string `json:"notes,omitempty" gorm:"type:text"`
}

func (osd *OrderShippingDetail) TableName() string {
	return "order_shipping_details"
}
