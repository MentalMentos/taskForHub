package repository

import (
	"context"
	"github.com/MentalMentos/taskForHub/books/internal/model"
)

type Books interface {
	Create(ctx context.Context, book *model.Book) error
	GetAll(ctx context.Context) ([]model.Book, error)
	GetByID(ctx context.Context, id string) (*model.Book, error)
	Delete(ctx context.Context, id string) error
}
