package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/MentalMentos/taskForHub/books/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository struct {
	DB *mongo.Database
}

func NewBookRepository(db *mongo.Database) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) Create(ctx context.Context, book *model.Book) error {
	collection := r.DB.Collection("books")

	book.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(ctx, book)
	if err != nil {
		return fmt.Errorf("failed to insert book: %w", err)
	}

	return nil
}

func (r *BookRepository) GetByID(ctx context.Context, bookID string) (*model.Book, error) {
	collection := r.DB.Collection("books")

	objectID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return nil, errors.New("invalid book ID format")
	}

	filter := bson.M{"_id": objectID}

	var book model.Book
	err = collection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("book not found")
		}
		return nil, fmt.Errorf("failed to get book by ID: %w", err)
	}

	return &book, nil
}

func (r *BookRepository) GetAll(ctx context.Context) ([]model.Book, error) {
	collection := r.DB.Collection("books")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}
	defer cursor.Close(ctx)

	var books []model.Book
	for cursor.Next(ctx) {
		var book model.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("failed to decode book: %w", err)
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return books, nil
}
