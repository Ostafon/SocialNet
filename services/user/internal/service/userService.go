package service

import (
	"bytes"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"socialnet/pkg/storage"
	"socialnet/pkg/utils"
	pb "socialnet/services/user/gen"
	"socialnet/services/user/internal/model"
	"socialnet/services/user/internal/repos"
)

type UserService struct {
	repo *repos.UserRepo
}

func NewUserService(r *repos.UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetUser(req *pb.GetUserRequest) (*model.User, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "required id")
	}
	id, err := utils.StringToUint(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(req *pb.UpdateUserRequest) error {
	err := utils.ValidateStruct(req)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUser(req)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateAvatar(req *pb.UpdateAvatarRequest) (*pb.UpdateAvatarResponse, error) {
	// avatar → io.Reader

	reader := bytes.NewReader(req.Avatar)
	s3, err := storage.NewS3Client()
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot to create client")
	}

	// имя файла
	key := fmt.Sprintf("avatars/%s_%s", req.Id, req.Filename)

	// загрузка в S3
	url, err := s3.UploadFile(reader, key)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload avatar: %v", err)
	}

	// обновляем в БД
	if err := s.repo.UpdateAvatar(req.Id, url); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update avatar in DB: %v", err)
	}

	return &pb.UpdateAvatarResponse{
		AvatarUrl: url,
	}, nil
}

func (s *UserService) FollowUser(followerID, followingID string) error {

	if followerID == followingID {
		return fmt.Errorf("cannot follow yourself")
	}

	follower, _ := utils.StringToUint(followerID)
	following, _ := utils.StringToUint(followingID)
	if err := s.repo.FollowByUserId(follower, following); err != nil {
		return err
	}
	return nil
}

func (s *UserService) UnfollowUser(followerID, followingID string) error {

	follower, _ := utils.StringToUint(followerID)
	following, _ := utils.StringToUint(followingID)
	if err := s.repo.DeleteFollow(follower, following); err != nil {
		return err
	}
	return nil
}
