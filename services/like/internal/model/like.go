package model

import (
	"time"
)

type Like struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index;not null"`
	PostID    *string   `gorm:"index"`
	CommentID *string   `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
