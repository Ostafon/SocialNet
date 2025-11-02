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
