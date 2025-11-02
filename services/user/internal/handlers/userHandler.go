package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb1 "socialnet/services/auth/gen"
	pb "socialnet/services/user/gen"
	"socialnet/services/user/internal/service"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	serv *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{serv: s}
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {

	user, err := h.serv.GetUser(req)
	if err != nil {
		return nil, err
	}
	return &pb.User{FirstName: user.Firstname, LastName: user.Lastname,
		BirthDate: user.BirthDate, Bio: user.Bio, AvatarUrl: user.AvatarUrl}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb1.Confirmation, error) {
	err := h.serv.UpdateUser(req)
	if err != nil {
		return nil, err
	}
	return &pb1.Confirmation{Status: codes.OK.String()}, nil
}

func (h *UserHandler) UpdateAvatar(ctx context.Context, req *pb.UpdateAvatarRequest) (*pb.UpdateAvatarResponse, error) {
	avatar, err := h.serv.UpdateAvatar(req)
	if err != nil {
		return nil, err
	}
	return avatar, nil
}

func (h *UserHandler) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb1.Confirmation, error) {
	userId := ctx.Value("user_id").(string)
	if userId == "" {
		return nil, status.Error(codes.Internal, "user id is null")
	}
	fmt.Println("GetProfile userId:", userId)
	err := h.serv.FollowUser(userId, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb1.Confirmation{Status: "Followed Successfully"}, nil
}

func (h *UserHandler) UnfollowUser(ctx context.Context, req *pb.UnfollowUserRequest) (*pb1.Confirmation, error) {
	userId := ctx.Value("user_id").(string)
	if userId == "" {
		return nil, status.Error(codes.Internal, "user id is null")
	}
	fmt.Println("GetProfile userId:", userId)
	if err := h.serv.UnfollowUser(userId, req.Id); err != nil {
		return nil, err
	}
	return &pb1.Confirmation{Status: "Successfully"}, nil
}
