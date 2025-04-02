package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/MentalMentos/taskForHub/auth/internal/model"
	"github.com/MentalMentos/taskForHub/auth/pkg/helpers"
	"github.com/MentalMentos/taskForHub/auth/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoImpl struct {
	DB     *mongo.Database
	logger logger.Logger
}

func NewRepo(db *mongo.Database, logger logger.Logger) *RepoImpl {
	return &RepoImpl{
		DB:     db,
		logger: logger,
	}
}

func (r *RepoImpl) Create(ctx context.Context, user model.User) (string, error) {
	collection := r.DB.Collection("users")

	// Создаем новый ObjectID для пользователя
	user.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		r.logger.Info(helpers.RepoPrefix, err.Error())
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	r.logger.Info("repo_create_user", "successful")
	return user.ID.Hex(), nil
}

func (r *RepoImpl) GetByEmail(ctx context.Context, email string) (model.User, string, error) {
	collection := r.DB.Collection("users")

	filter := bson.M{"email": email}

	var user model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, "", errors.New("user not found")
		}
		r.logger.Info("Failed to get user by email:", err.Error())
		return model.User{}, "", fmt.Errorf("cannot get user by email: %w", err)
	}
	userID := user.ID.Hex()
	return user, userID, nil
}

func (r *RepoImpl) GetByID(ctx context.Context, userID string) (model.User, error) {
	collection := r.DB.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.User{}, errors.New("invalid user ID format")
	}

	filter := bson.M{"_id": objectID}

	var user model.User
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, errors.New("user not found")
		}
		r.logger.Info("Failed to get user by ID:", err.Error())
		return model.User{}, fmt.Errorf("cannot get user by ID: %w", err)
	}

	return user, nil
}

func (r *RepoImpl) GetAll(ctx context.Context) ([]model.User, error) {
	collection := r.DB.Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		r.logger.Info("Failed to get all users:", err.Error())
		return nil, fmt.Errorf("cannot get all users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []model.User
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			r.logger.Info("Failed to decode user:", err.Error())
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Info("Cursor error:", err.Error())
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}
