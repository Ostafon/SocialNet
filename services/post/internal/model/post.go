package model

import "time"

type Post struct {
	ID            uint   `gorm:"primaryKey"`
	UserId        string `gorm:"index;not null"`
	Content       string `gorm:"type:text"`
	ImageUrl      string
	LikesCount    int32     `gorm:"default:0"`
	CommentsCount int32     `gorm:"default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
