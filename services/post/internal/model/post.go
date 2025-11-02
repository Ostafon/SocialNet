package model

type Post struct {
	Id            uint   `gorm:"primaryKey"`
	UserId        string `gorm:"index;not null"`
	Content       string `gorm:"type:text"`
	ImageUrl      string
	LikesCount    int32
	CommentsCount int32
	CreatedAt     string
	UpdatedAt     string
}
