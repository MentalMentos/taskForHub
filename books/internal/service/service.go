package service

import (
	"context"
	"github.com/MentalMentos/taskForHub/books/internal/model"
)

type Book interface {
	CreateBook(ctx context.Context, book *model.Book) error
	GetAllBooks(ctx context.Context) ([]model.Book, error)
	GetBookByID(ctx context.Context, id string) (*model.Book, error)
}
