package repository

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/MentalMentos/taskForHub/auth/internal/model"
	"github.com/MentalMentos/taskForHub/auth/pkg/helpers"
	"github.com/MentalMentos/taskForHub/auth/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RepoImpl struct {
	DB     *mongo.Database
	logger logger.Logger
}

func NewRepo(db *mongo.Database, logger logger.Logger) *RepoImpl {
	return &RepoImpl{
		db,
		logger,
	}
}

func (r *RepoImpl) Create(ctx context.Context, user model.User) (int64, error) {
	collection := r.DB.Collection("users")
	insertResult, err := collection.InsertOne(ctx, bson.D{
		{Key: "name", Value: user.Name},
		{Key: "email", Value: user.Email},
		{Key: "password", Value: user.Password},
	})
	if err != nil {
		r.logger.Info(helpers.RepoPrefix, err.Error())
		return 0, err
	}
	objectID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		r.logger.Info(helpers.RepoPrefix, "Inserted ID is not primitive.ObjectID")
		return 0, errors.New("cannot convert InsertedID to ObjectID")
	}

	objectIDBytes := objectID[:]

	id := int64(binary.BigEndian.Uint64(objectIDBytes[4:]))
	r.logger.Info("repo_create_user", "successful")
	return id, nil
}

func (r *RepoImpl) GetByEmail(ctx context.Context, email string) (model.User, error) {
	collection := r.DB.Collection("users")

	filter := bson.M{"email": email}

	var user model.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, errors.New("user not found")
		}
		r.logger.Info("Failed to get user by email: %v", err.Error())
		return model.User{}, fmt.Errorf("cannot get user by email: %w", err)
	}

	return user, nil
}

func (r *RepoImpl) GetByID(ctx context.Context, userID int64) (model.User, error) {
	collection := r.DB.Collection("users")

	objectID := primitive.NewObjectIDFromTimestamp(time.Unix(userID, 0))

	filter := bson.M{"_id": objectID}

	var user model.User
	err := collection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, errors.New("user not found")
		}
		r.logger.Info("Failed to get user by ID: %v", err.Error())
		return model.User{}, fmt.Errorf("cannot get user by ID: %w", err)
	}

	return user, nil
}

func (r *RepoImpl) GetAll(ctx context.Context) ([]model.User, error) {
	collection := r.DB.Collection("users")

	findOptions := options.Find()
	var users []model.User

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		r.logger.Info("Failed to get all users: %v", err.Error())
		return nil, fmt.Errorf("cannot get all users: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model.User
		err := cursor.Decode(&user)
		if err != nil {
			r.logger.Info("Failed to decode user: %v", err.Error())
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Info("Cursor error: %v", err.Error())
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}
