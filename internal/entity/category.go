package entity

import "time"

type Category struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	ParentID    *string   `json:"parent_id" gorm:"type:uuid;default:null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
}

func (c *Category) TableName() string {
	return "categories"
}
