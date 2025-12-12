package repos

import (
	"gorm.io/gorm"
	"socialnet/services/notification/internal/model"
)

type NotificationRepo struct {
	DB *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) *NotificationRepo {
	return &NotificationRepo{DB: db}
}

func (r *NotificationRepo) Save(n *model.Notification) error {
	return r.DB.Create(n).Error
}

func (r *NotificationRepo) List(userID, filter string, limit, offset int) ([]model.Notification, error) {
	query := r.DB.Where("user_id = ?", userID)

	if filter == "unread" {
		query = query.Where("read = ?", false)
	}

	if limit == 0 {
		limit = -1
	}

	var notifications []model.Notification
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	return notifications, err
}

func (r *NotificationRepo) MarkAsRead(id string) error {
	return r.DB.Model(&model.Notification{}).Where("id = ?", id).Update("read", true).Error
}

func (r *NotificationRepo) MarkAllAsRead(userID string) error {
	return r.DB.Model(&model.Notification{}).Where("user_id = ?", userID).Update("read", true).Error
}

func (r *NotificationRepo) Delete(id string) error {
	return r.DB.Delete(&model.Notification{}, id).Error
}

func (r *NotificationRepo) ClearAll(userID string) error {
	return r.DB.Where("user_id = ?", userID).Delete(&model.Notification{}).Error
}
