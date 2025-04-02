package repository

import (
	"context"
	"github.com/MentalMentos/taskForHub/books/internal/model"
	"github.com/kamva/mgm/v3"
)

type BookRepository struct{}

func NewBookRepository() *BookRepository {
	return &BookRepository{}
}

// Добавление книги
func (r *BookRepository) Create(ctx context.Context, book *model.Book) error {
	return mgm.Coll(book).CreateWithCtx(ctx, book)
}

// Получение всех книг
func (r *BookRepository) GetAll(ctx context.Context) ([]model.Book, error) {
	var books []model.Book
	err := mgm.Coll(&model.Book{}).SimpleFindWithCtx(ctx, &books, nil)
	return books, err
}

// Получение книги по ID
func (r *BookRepository) GetByID(ctx context.Context, id string) (*model.Book, error) {
	book := &model.Book{}
	err := mgm.Coll(book).FindByIDWithCtx(ctx, id, book)
	return book, err
}

// Удаление книги
func (r *BookRepository) Delete(ctx context.Context, id string) error {
	book, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return mgm.Coll(book).DeleteWithCtx(ctx, book)
}
