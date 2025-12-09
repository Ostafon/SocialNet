package repos

import (
	"gorm.io/gorm"
	"socialnet/services/like/internal/model"
)

type LikeRepo struct {
	db *gorm.DB
}

func NewLikeRepo(db *gorm.DB) *LikeRepo {
	return &LikeRepo{db: db}
}

// ---- POST ----
func (r *LikeRepo) LikePost(userID, postID string) error {
	like := &model.Like{UserID: userID, PostID: &postID}
	return r.db.Create(like).Error
}

func (r *LikeRepo) UnlikePost(userID, postID string) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Like{}).Error
}

func (r *LikeRepo) CountPostLikes(postID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

func (r *LikeRepo) ListPostLikes(postID string) ([]model.Like, error) {
	var likes []model.Like
	err := r.db.Where("post_id = ?", postID).Find(&likes).Error
	return likes, err
}
func (r *LikeRepo) HasUserLikedPost(userID, postID string) (bool, error) {
	var count int64
	err := r.db.
		Model(&model.Like{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ---- COMMENT ----
func (r *LikeRepo) LikeComment(userID, commentID string) error {
	like := &model.Like{UserID: userID, CommentID: &commentID}
	return r.db.Create(like).Error
}

func (r *LikeRepo) UnlikeComment(userID, commentID string) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&model.Like{}).Error
}

func (r *LikeRepo) CountCommentLikes(commentID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Like{}).Where("comment_id = ?", commentID).Count(&count).Error
	return count, err
}

func (r *LikeRepo) ListCommentLikes(commentID string) ([]model.Like, error) {
	var likes []model.Like
	err := r.db.Where("comment_id = ?", commentID).Find(&likes).Error
	return likes, err
}
