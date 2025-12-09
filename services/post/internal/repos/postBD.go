package repos

import (
	"gorm.io/gorm"
	"socialnet/services/post/internal/model"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) SavePost(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepo) GetPostByID(id string) (*model.Post, error) {
	post := &model.Post{}
	if err := r.db.Where("id = ?", id).First(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepo) GetAllPosts() ([]*model.Post, error) {
	var posts []*model.Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepo) DeletePostByID(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Post{}).Error
}

func (r *PostRepo) GetUserPosts(userID string) ([]*model.Post, error) {
	var posts []*model.Post
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepo) GetPostsByUsers(userIDs []string) ([]*model.Post, error) {
	var posts []*model.Post
	if err := r.db.
		Where("user_id IN ?", userIDs).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepo) SearchPosts(query string, limit, offset int) ([]model.Post, error) {
	var posts []model.Post
	q := "%" + query + "%"
	err := r.db.Where(`
		content ILIKE ? OR
		title ILIKE ?`, q, q,
	).Limit(limit).Offset(offset).Find(&posts).Error
	return posts, err
}
