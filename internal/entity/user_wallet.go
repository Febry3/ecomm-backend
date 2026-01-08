package entity

import "time"

type UserWallet struct {
	UserID    int64     `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	Balance   float64   `json:"balance" gorm:"type:decimal(15,2);default:0;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (w *UserWallet) TableName() string {
	return "user_wallets"
}
