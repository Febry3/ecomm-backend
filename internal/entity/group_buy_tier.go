package entity

type GroupBuyTier struct {
	ID                   string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GroupBuySessionID    string  `json:"group_buy_session_id" gorm:"type:uuid;not null"`
	ParticipantThreshold int     `json:"participant_threshold" gorm:"not null"` // e.g., 10, 50
	DiscountPercentage   float64 `json:"discount_percentage" gorm:"type:decimal(5,2);not null"`
}
