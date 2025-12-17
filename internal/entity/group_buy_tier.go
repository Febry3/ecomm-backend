package entity

type GroupBuyTier struct {
	ID                   string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GroupBuySessionID    string  `gorm:"type:uuid;not null"`
	ParticipantThreshold int     `gorm:"not null"` // e.g., 10, 50
	DiscountPercentage   float64 `gorm:"type:decimal(5,2);not null"`
}
