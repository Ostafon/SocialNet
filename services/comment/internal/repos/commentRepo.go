package repos

import (
	"gorm.io/gorm"
	"socialnet/services/comment/internal/model"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) AddComment(c *model.Comment) error {
	return r.db.Create(c).Error
}

func (r *CommentRepo) GetComment(id string) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepo) DeleteComment(id, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Comment{}).Error
}

func (r *CommentRepo) ListComments(postID string) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Where("post_id = ?", postID).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *CommentRepo) UpdateLikesCount(commentID string, count int) error {
	return r.db.Model(&model.Comment{}).
		Where("id = ?", commentID).
		Update("likes_count", count).
		Error
}
