package model

import (
	"time"
)

type Comment struct {
	ID         uint      `gorm:"primaryKey"`
	PostID     string    `gorm:"index;not null"`
	UserID     string    `gorm:"index;not null"`
	Content    string    `gorm:"type:text;not null"`
	LikesCount int       `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
