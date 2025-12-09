package model

import "time"

type Notification struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      string `gorm:"index"`
	Type        string
	ReferenceID string
	Content     string
	Read        bool `gorm:"default:false"`
	CreatedAt   time.Time
}
