package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;size:100;not null"`
	Username  string `gorm:"uniqueIndex;size:50;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
}
