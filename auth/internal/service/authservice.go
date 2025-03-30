package service

import (
	"context"
	"errors"
	"github.com/MentalMentos/taskForHub/auth/internal/data/request"
	"github.com/MentalMentos/taskForHub/auth/internal/data/response"
	"github.com/MentalMentos/taskForHub/auth/internal/model"
	"github.com/MentalMentos/taskForHub/auth/internal/repository"
	"github.com/MentalMentos/taskForHub/auth/pkg/helpers"
	"github.com/MentalMentos/taskForHub/auth/pkg/logger"
	"github.com/MentalMentos/taskForHub/auth/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Register(ctx context.Context, req request.RegisterUserRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req request.LoginRequest) (*model.AuthResponse, error)
	GetAccessToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error)
	UpdatePassword(ctx context.Context, req request.UpdateUserRequest) (*response.UpdatePasswordResponse, error)
}

type AuthService struct {
	repo   repository.Repository
	logger logger.Logger
}

func NewAuthService(repo repository.Repository, logger logger.Logger) *AuthService {
	return &AuthService{
		repo,
		logger,
	}
}

func (s *AuthService) Register(ctx context.Context, req request.RegisterUserRequest) (*model.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Fatal("[ SERVICE_REGISTER ]", helpers.FailedToHashPass)
		return nil, err
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
		IP:       "192.168.123.132",
	}

	if req.Role == "" {
		req.Role = "user"
	}

	_, err = s.repo.Create(ctx, user, s.logger)
	if err != nil {
		s.logger.Fatal("[ SERVICE_REGISTER ]", helpers.FailedToCreateUser)
		return nil, err
	}

	accessToken, refreshToken, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		s.logger.Fatal("[ SERVICE_REGISTER ]", helpers.FailedToGenJWT)
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email, s.logger)
	if err != nil {
		s.logger.Fatal("[ SERVICE_LOGIN ]", helpers.FailedToGetUser)
		return nil, errors.New("user not found")
	}

	if user.IP != req.IP {
		_, err := s.repo.UpdateIP(ctx, user, req.IP, s.logger)
		if err != nil {
			s.logger.Fatal("[ SERVICE_LOGIN ]", "failed to update ip")
			return nil, errors.New("cannot update ip with login")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.Fatal("[ SERVICE_LOGIN ]", helpers.FailedToHashPass)
		return nil, errors.New("invalid password")
	}

	accessToken, refreshToken, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		s.logger.Fatal("[ SERVICE_LOGIN ]", helpers.FailedToGenJWT)
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, req request.UpdateUserRequest) (*response.UpdatePasswordResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email, s.logger)
	if err != nil {
		s.logger.Fatal("[ SERVICE_UPDATE_PASSWORD ]", helpers.FailedToGetUser)
		return nil, errors.New("user not found")
	}

	if user.IP != req.IP {
		_, err := s.repo.UpdateIP(ctx, user, req.IP, s.logger)
		if err != nil {
			s.logger.Fatal("[ SERVICE_UPDATE_PASSWORD ]", "failed to update ip with login")
			return nil, errors.New("cannot update ip with login")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.Fatal("[ SERVICE_UPDATE_PASSWORD ]", "invalid password")
		return nil, errors.New("invalid password")
	}
	accessToken, refreshToken, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		s.logger.Fatal("[ SERVICE_UPDATE_PASSWORD ]", "failed to generate access token")
		return nil, err
	}
	return &response.UpdatePasswordResponse{
		accessToken,
		refreshToken,
		req.Name,
	}, nil
}

// Метод для обновления access token
func (s *AuthService) GetAccessToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error) {
	// Валидация refresh token
	claims, err := utils.ValidateJWT(refreshToken)
	if err != nil {
		s.logger.Fatal("[ SERVICE_GET_ACCESS_TOKEN ]", "failed to validate tokens")
		return nil, errors.New("invalid refresh token")
	}

	// Генерация нового набора токенов
	newAccessToken, newRefreshToken, err := utils.GenerateJWT(claims.UserID, claims.Role)
	if err != nil {
		s.logger.Fatal("[ SERVICE_GET_ACCESS_TOKEN ]", "failed to generate access tokens")
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
