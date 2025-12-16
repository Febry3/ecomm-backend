package entity

import "time"

type Category struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	ParentID    *int64    `json:"parent_id" gorm:"default:null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (c *Category) TableName() string {
	return "categories"
}
