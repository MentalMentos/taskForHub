package service

import (
	"context"
	"github.com/MentalMentos/taskForHub/auth/internal/data/request"
	"github.com/MentalMentos/taskForHub/auth/internal/data/response"
	"github.com/MentalMentos/taskForHub/auth/internal/model"
	"github.com/MentalMentos/taskForHub/auth/internal/repository"
	"github.com/MentalMentos/taskForHub/auth/pkg/logger"
)

type Auth interface {
	Register(ctx context.Context, req request.RegisterUserRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req request.LoginRequest) (*model.AuthResponse, error)
	GetAccessToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error)
}

type Service struct {
	*AuthService
}

func New(repo repository.Repository, logger logger.Logger) *Service {
	return &Service{
		NewAuthService(repo, logger),
	}
}
