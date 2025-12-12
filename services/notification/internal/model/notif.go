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

type RedisNotification struct {
	Id          string `json:"id"`
	UserID      string `json:"user_id"`
	Type        string `json:"type"`
	ReferenceID string `json:"reference_id"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
}
