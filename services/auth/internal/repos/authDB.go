package repos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"socialnet/pkg/utils"
	"socialnet/services/auth/internal/model"
	"time"
)

type UserRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) RegisterDB(user *model.User) (uint, error) {

	err := r.db.Create(user).Error
	if err != nil {
		return 0, utils.ErrorHandler(err, "Failed to register")
	}

	return user.ID, nil
}

func (r *UserRepo) SaveToken(token string, id uint) error {
	// Удаляем старые токены пользователя
	if err := r.db.Where("user_id = ?", id).Delete(&model.RefreshToken{}).Error; err != nil {
		return utils.ErrorHandler(err, "Cannot clear old refresh tokens")
	}

	refresh := &model.RefreshToken{
		Token:     token,
		UserID:    id,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := r.db.Create(refresh).Error; err != nil {
		return utils.ErrorHandler(err, "Cannot add refresh token")
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetUserById(id uint) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) FindRefreshTokenById(id uint) (*model.RefreshToken, error) {
	refresh := &model.RefreshToken{}
	if err := r.db.Where("user_id = ?", id).First(refresh).Error; err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return refresh, nil
}

func (r *UserRepo) CheckRefreshToken(refresh string) (uint, error) {
	token := &model.RefreshToken{}

	// Сначала ищем в БД
	if err := r.db.Where("token = ?", refresh).First(token).Error; err != nil {
		return 0, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	// Потом проверяем срок действия
	if time.Now().After(token.ExpiresAt) {
		return 0, status.Error(codes.Unauthenticated, "refresh token expired")
	}

	return token.UserID, nil
}

func (r *UserRepo) UpdatePassword(id uint, newPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPassword).Error
}

func (r *UserRepo) SaveResetToken(userID uint, token string) error {
	reset := &model.PasswordReset{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	return r.db.Create(reset).Error
}

func (r *UserRepo) FindResetToken(token string) (*model.PasswordReset, error) {
	reset := &model.PasswordReset{}
	if err := r.db.Where("token = ?", token).First(reset).Error; err != nil {
		return nil, err
	}
	return reset, nil
}

func (r *UserRepo) DeleteResetToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&model.PasswordReset{}).Error
}
