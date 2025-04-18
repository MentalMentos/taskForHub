package repository

import (
	"context"
	"github.com/MentalMentos/taskForHub/auth/internal/model"
	"github.com/MentalMentos/taskForHub/auth/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, us model.User) (string, error)
	GetByEmail(ctx context.Context, email string) (model.User, string, error)
	GetByID(ctx context.Context, userID string) (model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
}

type Repo struct {
	Repository
}

func NewRepository(db *mongo.Database, mylogger logger.Logger) *Repo {
	return &Repo{NewRepo(db, mylogger)}
}
