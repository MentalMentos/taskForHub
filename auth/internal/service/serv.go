package service

import (
	"task0325/auth/internal/repository"
	"task0325/auth/pkg/logger"
)

type Service struct {
	*AuthService
}

func New(repo repository.Repository, logger logger.Logger) *Service {
	return &Service{
		NewAuthService(repo, logger),
	}
}
