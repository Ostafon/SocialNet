package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"socialnet/pkg/contextx"
	"socialnet/pkg/utils"
	pb "socialnet/services/auth/gen"
	"socialnet/services/auth/internal/service"
	"time"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	fmt.Println("Register")
	access, refresh, err := h.authService.Register(req)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{RefreshToken: refresh, AccessToken: access}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Println("Login")
	access, refresh, err := h.authService.Login(req)
	if err != nil {
		return nil, err
	}
	return &pb.LoginResponse{RefreshToken: refresh, AccessToken: access}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	token, err := h.authService.RefreshToken(req)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	resp := &pb.RefreshResponse{AccessToken: token}
	return resp, nil
}

func (h *AuthHandler) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.UpdatePasswordResponse, error) {
	token, err := h.authService.UpdatePassword(req)
	if err != nil {
		return nil, err
	}
	resp := &pb.UpdatePasswordResponse{AccessToken: token, Message: "successfully"}
	return resp, nil
}

func (h *AuthHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	err := h.authService.ForgotPassword(req)
	if err != nil {
		return nil, err
	}

	return &pb.ForgotPasswordResponse{Status: status.New(codes.OK, "password reset successful").String()}, nil
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.Confirmation, error) {
	err := h.authService.ResetPassword(req)
	fmt.Println("Reset")
	if err != nil {
		return nil, err
	}
	return &pb.Confirmation{Status: status.New(codes.OK, "password updated successful").String()}, nil
}

func (h *AuthHandler) GetProfile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	userId := contextx.GetUserID(ctx)
	if userId == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user_id in context")
	}

	fmt.Println("GetProfile userId:", userId)

	id, err := utils.StringToUint(userId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	user, err := h.authService.GetProfile(id)
	if err != nil {
		return nil, err
	}

	return &pb.ProfileResponse{
		Id:        userId,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
