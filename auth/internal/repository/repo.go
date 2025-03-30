package repository

import (
	"context"
	"gorm.io/gorm"
	"task0325/auth/internal/model"
	"task0325/auth/pkg/logger"
)

type Repository interface {
	Create(ctx context.Context, us model.User, logger logger.Logger) (int64, error)
	Update(ctx context.Context, us model.User, logger logger.Logger) (int64, error)
	Delete(ctx context.Context, usId int64, logger logger.Logger) error
	UpdatePassword(ctx context.Context, us model.User, hashPassword string, logger logger.Logger) (model.User, error)
	UpdateIP(ctx context.Context, us model.User, ip string, logger logger.Logger) (model.User, error)
	GetByEmail(ctx context.Context, email string, logger logger.Logger) (model.User, error)
	GetByID(ctx context.Context, userID int64, logger logger.Logger) (model.User, error)
	GetAll(ctx context.Context, logger logger.Logger) ([]model.User, error)
}

type Repo struct {
	Repository
}

func NewRepository(db *gorm.DB, mylogger logger.Logger) *Repo {
	return &Repo{NewRepo(db, mylogger)}
}
