package repository

import (
	"context"
	"fmt"
	"github.com/MentalMentos/taskForHub/cart/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CartRepository - структура репозитория корзины
type CartRepository struct {
	DB *mongo.Database
}

// NewCartRepository - конструктор для CartRepository
func NewCartRepository(db *mongo.Database) *CartRepository {
	return &CartRepository{DB: db}
}

// AddToCart добавляет товар в корзину. Если корзины нет, она создается.
func (r *CartRepository) AddToCart(ctx context.Context, userID primitive.ObjectID, item model.CartItem) error {
	// Находим корзину по user_id
	collection := r.DB.Collection("carts")

	var cart model.Cart
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		// Если корзины нет, создаем новую корзину с этим товаром
		cart = model.Cart{
			Items: []model.CartItem{item},
		}
		_, err := collection.InsertOne(ctx, bson.M{"_id": userID, "items": cart.Items})
		if err != nil {
			return fmt.Errorf("failed to insert cart: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to find cart: %w", err)
	}

	// Если корзина существует, добавляем товар в существующую корзину
	cart.Items = append(cart.Items, item)

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{"items": cart.Items},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}
	return nil
}

// GetCart возвращает корзину пользователя по его user_id
func (r *CartRepository) GetCart(ctx context.Context, userID primitive.ObjectID) (*model.Cart, error) {
	collection := r.DB.Collection("carts")

	var cart model.Cart
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("cart not found for user %s", userID.Hex())
	} else if err != nil {
		return nil, fmt.Errorf("failed to find cart: %w", err)
	}

	return &cart, nil
}
