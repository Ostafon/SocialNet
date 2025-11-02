package repos

import (
	"gorm.io/gorm"
	pb "socialnet/services/post/gen"
	"socialnet/services/post/internal/model"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) SavePost(post *model.Post) error {

}

func (r *PostRepo) GetPostById(id uint) (*pb.Post, error) {

}
