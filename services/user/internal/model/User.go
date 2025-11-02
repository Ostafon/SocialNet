package model

import "time"

type User struct {
	Id        uint   `gorm:"primaryKey"`
	Firstname string `gorm:"size:50;not null"`
	Lastname  string `gorm:"size:50;not null"`
	BirthDate string `gorm:"size:50"`
	Bio       string `gorm:"size:255;"`
	AvatarUrl string `gorm:"size:255"`
	CreatedAt time.Time
}
