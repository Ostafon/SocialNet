package repos

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"socialnet/pkg/utils"
	pb "socialnet/services/user/gen"
	"socialnet/services/user/internal/model"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUser(id uint) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return user, nil
}

func (r *UserRepo) UpdateUser(req *pb.UpdateUserRequest) error {
	id, err := utils.StringToUint(req.Id)
	if err != nil {
		return err
	}
	user := &model.User{Id: id, Firstname: req.FirstName, Lastname: req.LastName,
		BirthDate: req.BirthDate, Bio: req.Bio}
	if err := r.db.Where("id = ?", req.Id).Save(&user).Error; err != nil {
		return status.Error(codes.Internal, "unable to save info")
	}
	return nil
}

func (r *UserRepo) UpdateAvatar(userID string, avatarURL string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error
}

func (r *UserRepo) FollowByUserId(followerID, followingID uint) error {
	follow := &model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := r.db.Create(follow).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) DeleteFollow(followerID, followingID uint) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&model.Follow{}).Error
}

func (r *UserRepo) GetFollowing(userID string) ([]*model.User, error) {
	var users []*model.User

	err := r.db.
		Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
func (r *UserRepo) SearchUsers(query string, limit, offset int) ([]model.User, error) {
	var users []model.User
	q := "%" + query + "%"
	err := r.db.Where(`
		username ILIKE ? OR
		first_name ILIKE ? OR
		last_name ILIKE ? OR
		CONCAT(first_name, ' ', last_name) ILIKE ? OR
		CONCAT(last_name, ' ', first_name) ILIKE ?`,
		q, q, q, q, q,
	).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *UserRepo) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
