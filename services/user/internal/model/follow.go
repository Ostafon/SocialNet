package model

import "time"

type Follow struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	FollowerID  uint      `gorm:"not null,index:idx_follows_unique,unique"`
	FollowingID uint      `gorm:"not null,index:idx_follows_unique,unique"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
