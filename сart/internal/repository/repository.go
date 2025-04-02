package repository

import (
	"context"
	"errors"
	"github.com/MentalMentos/taskForHub/сart/internal/model"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

func (r *CartRepository) GetCart(ctx context.Context, userID primitive.ObjectID) (*model.Cart, error) {
	var cart model.Cart
	err := mgm.Coll(&cart).First(bson.M{"user_id": userID}, &cart)
	if err != nil {
		return nil, errors.New("cart not found")
	}
	return &cart, nil
}

func (r *CartRepository) AddToCart(ctx context.Context, userID primitive.ObjectID, item model.CartItem) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		cart = &model.Cart{
			UserID: userID,
			Items:  []model.CartItem{item},
		}
		return mgm.Coll(cart).Create(cart)
	}
	cart.Items = append(cart.Items, item)
	return mgm.Coll(cart).Update(cart)
}

func (r *CartRepository) RemoveFromCart(ctx context.Context, userID primitive.ObjectID, productID primitive.ObjectID) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	var updatedItems []model.CartItem
	for _, item := range cart.Items {
		if item.ProductID != productID {
			updatedItems = append(updatedItems, item)
		}
	}
	cart.Items = updatedItems
	return mgm.Coll(cart).Update(cart)
}
