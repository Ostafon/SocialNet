package model

import "time"

type Chat struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100"`
	IsGroup      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Participants []Participant `gorm:"foreignKey:ChatID"`
}

type Participant struct {
	ID     uint `gorm:"primaryKey"`
	ChatID uint
	UserID string `gorm:"index"`
}

type Message struct {
	ID          uint   `gorm:"primaryKey"`
	ChatID      uint   `gorm:"index"`
	SenderID    string `gorm:"index"`
	Content     string
	ContentType string
	MediaURL    string
	Read        bool
	CreatedAt   time.Time
}
