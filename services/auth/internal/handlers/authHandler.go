package handlers

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "socialnet/services/auth/gen"
	"socialnet/services/auth/internal/service"
	"socialnet/services/auth/internal/utils"
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
	access, refresh, err := h.authService.Register(req)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterResponse{RefreshToken: refresh, AccessToken: access}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
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
	if err != nil {
		return nil, err
	}
	return &pb.Confirmation{Status: status.New(codes.OK, "password updated successful").String()}, nil
}

func (h *AuthHandler) GetProfile(ctx context.Context, req *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	id, err := utils.StringToUint(userID)
	user, err := h.authService.GetProfile(id)
	if err != nil {
		return nil, err
	}

	return &pb.ProfileResponse{
		Id:        fmt.Sprint(user.ID),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
