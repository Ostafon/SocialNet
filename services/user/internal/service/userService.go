package service

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"socialnet/pkg/config"
	"socialnet/pkg/storage"
	"socialnet/pkg/utils"
	notificationpb "socialnet/services/notification/gen"
	pb "socialnet/services/user/gen"
	"socialnet/services/user/internal/model"
	"socialnet/services/user/internal/repos"
)

type UserService struct {
	repo    *repos.UserRepo
	clients *config.GRPCClients
}

func NewUserService(r *repos.UserRepo, clients *config.GRPCClients) *UserService {
	return &UserService{repo: r, clients: clients}
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

func (s *UserService) FollowUser(ctx context.Context, followerID, followingID string) error {

	if followerID == followingID {
		return fmt.Errorf("cannot follow yourself")
	}

	follower, _ := utils.StringToUint(followerID)
	following, _ := utils.StringToUint(followingID)
	if err := s.repo.FollowByUserId(follower, following); err != nil {
		return err
	}

	md := metadata.New(map[string]string{"user-id": followerID})
	ctxWithUser := metadata.NewOutgoingContext(ctx, md)

	notifClient, err := s.clients.GetNotifClient("localhost:50057")
	if err == nil {
		_, _ = notifClient.CreateNotification(ctxWithUser,
			&notificationpb.CreateNotificationRequest{
				UserId:      followingID,
				Type:        "follow",
				ReferenceId: followerID,
				Content:     fmt.Sprintf("User %s followed you", followerID),
			})
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

func (s *UserService) GetFollowing(userID string) ([]*pb.User, error) {
	users, err := s.repo.GetFollowing(userID)
	if err != nil {
		return nil, err
	}

	// Преобразуем в protobuf
	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id: fmt.Sprint(u.Id),
		})
	}
	return pbUsers, nil
}

func (s *UserService) GetAllUsers() (*pb.Users, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var pbUsers []*pb.User
	for _, v := range users {
		pbUsers = append(
			pbUsers,
			&pb.User{
				Id:        fmt.Sprint(v.Id),
				FirstName: v.Firstname,
				LastName:  v.Lastname,
				BirthDate: v.BirthDate,
				Bio:       v.Bio,
			},
		)
	}

	return &pb.Users{Users: pbUsers}, nil
}
