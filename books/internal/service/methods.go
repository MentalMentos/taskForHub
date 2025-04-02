package service

import (
	"context"
	"github.com/MentalMentos/taskForHub/books/internal/model"
	"github.com/MentalMentos/taskForHub/books/internal/repository"
	"github.com/labstack/gommon/log"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(ctx context.Context, book *model.Book) error {
	log.Info("service", book)
	return s.repo.Create(ctx, book)
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]model.Book, error) {
	return s.repo.GetAll(ctx)
}

func (s *BookService) GetBookByID(ctx context.Context, id string) (*model.Book, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BookService) DeleteBook(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
