package service

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	utils2 "socialnet/pkg/utils"
	pb "socialnet/services/auth/gen"
	"socialnet/services/auth/internal/model"
	"socialnet/services/auth/internal/repos"
	"socialnet/services/auth/internal/utils"
	"time"
)

type AuthService struct {
	repo *repos.UserRepo
}

func NewAuthService(repo *repos.UserRepo) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(req *pb.RegisterRequest) (string, string, error) {
	// валидация
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !re.MatchString(req.Email) {
		return "", "", status.Error(codes.InvalidArgument, "invalid mail")
	}
	if err := utils.ValidateStruct(req); err != nil {
		return "", "", status.Error(codes.InvalidArgument, "required all fields")
	}

	// хэшируем пароль
	encPass, err := utils.PasswordHashing(req.Password)
	if err != nil {
		return "", "", status.Error(codes.Internal, "internal error")
	}

	// создаём модель
	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: encPass,
	}

	// сохраняем в БД
	userID, err := s.repo.RegisterDB(user)
	if err != nil {
		return "", "", status.Error(codes.Internal, "db error")
	}

	// создаём access token
	accessToken, err := utils2.SignToken(fmt.Sprint(userID), user.Username)
	if err != nil {
		return "", "", status.Error(codes.Internal, "internal error")
	}

	// создаём refresh token
	refreshToken, err := utils2.GenerateRefreshToken()
	if err != nil {
		return "", "", status.Error(codes.Internal, "internal error")
	}

	// сохраняем refresh в БД
	if err = s.repo.SaveToken(refreshToken, userID); err != nil {
		return "", "", status.Error(codes.Internal, "db error")
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Login(req *pb.LoginRequest) (string, string, error) {
	// валидация
	if err := utils.ValidateStruct(req); err != nil {
		return "", "", status.Error(codes.InvalidArgument, "required all fields")
	}
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", "", status.Error(codes.NotFound, "user not exists")
	}

	if err = utils.VerifyPassword(user.Password, req.Password); err != nil {
		return "", "", status.Error(codes.InvalidArgument, "email or password wrong")
	}

	refresh, err := s.repo.FindRefreshTokenById(user.ID)
	if err != nil {
		return "", "", err
	}
	var refreshToken string
	if time.Now().Unix() > refresh.ExpiresAt.Unix() {
		refreshToken, err = utils2.GenerateRefreshToken()
		if err != nil {
			return "", "", status.Error(codes.Internal, "internal error")
		}
		if err := s.repo.SaveToken(refreshToken, user.ID); err != nil {
			return "", "", err
		}
	} else {
		refreshToken = refresh.Token
	}

	accessToken, err := utils2.SignToken(fmt.Sprint(user.ID), user.Username)
	if err != nil {
		return "", "", status.Error(codes.Internal, "internal error")
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(req *pb.RefreshRequest) (string, error) {
	id, err := s.repo.CheckRefreshToken(req.RefreshToken)
	if err != nil {
		return "", utils.ErrorHandler(err, "invalid token ")
	}

	user, err := s.repo.GetUserById(id)
	if err != nil {
		return "", status.Error(codes.NotFound, "user not found")
	}

	token, err := utils2.SignToken(fmt.Sprint(user.ID), user.Username)
	if err != nil {
		return "", utils.ErrorHandler(err, "internal error")
	}

	return token, nil
}

func (s *AuthService) UpdatePassword(req *pb.UpdatePasswordRequest) (string, error) {
	uid, err := utils.StringToUint(req.Id)
	if err != nil {
		return "", status.Error(codes.Internal, "internal error")
	}
	user, err := s.repo.GetUserById(uid)
	if err != nil {
		return "", status.Error(codes.NotFound, "user not found")
	}

	if err := utils.VerifyPassword(user.Password, req.CurrentPassword); err != nil {
		return "", status.Error(codes.InvalidArgument, "current password incorrect")
	}

	hashed, err := utils.PasswordHashing(req.NewPassword)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to hash password")
	}

	if err := s.repo.UpdatePassword(user.ID, hashed); err != nil {
		return "", status.Error(codes.Internal, "failed to update password")
	}

	access, err := utils2.SignToken(fmt.Sprint(user.ID), user.Username)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to generate token")
	}
	return access, nil
}

func (s *AuthService) ForgotPassword(req *pb.ForgotPasswordRequest) error {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return status.Error(codes.NotFound, "user not found")
	}

	token, err := utils.GenerateUUID()
	resetLink := "http://localhost:8080/reset-password?token=" + token
	if err != nil {
		return status.Error(codes.Internal, "failed to generate reset token")
	}
	if err := s.repo.SaveResetToken(user.ID, token); err != nil {
		return status.Error(codes.Internal, "failed to save reset token")
	}

	// Отправляем письмо
	subject := "Password Reset Request"
	body := "<p>Для сброса пароля перейдите по ссылке:</p>" +
		"<a href='" + resetLink + "'>" + resetLink + "</a>"

	if err = utils.SendEmail("sashaostafi60@gmail.com", subject, body); err != nil {
		return status.Error(codes.Internal, "failed to send email")
	}

	return nil
}

func (s *AuthService) ResetPassword(req *pb.ResetPasswordRequest) error {
	reset, err := s.repo.FindResetToken(req.ResetToken)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid reset token")
	}

	if reset.ExpiresAt.Before(time.Now()) {
		return status.Error(codes.InvalidArgument, "reset token expired")
	}

	hashed, err := utils.PasswordHashing(req.NewPassword)
	if err != nil {
		return status.Error(codes.Internal, "failed to hash password")
	}

	if err := s.repo.UpdatePassword(reset.UserID, hashed); err != nil {
		return status.Error(codes.Internal, "failed to update password")
	}

	if err := s.repo.DeleteResetToken(reset.Token); err != nil {
		return status.Error(codes.Internal, "failed to remove reset token")
	}

	return nil
}

func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return user, nil
}
