package service

import (
	"context"
	"github.com/MentalMentos/taskForHub/cart/internal/model"
	"github.com/MentalMentos/taskForHub/cart/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService struct {
	repo *repository.CartRepository
}

func NewCartService(repo *repository.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddToCart(ctx context.Context, userID primitive.ObjectID, item model.CartItem) error {
	return s.repo.AddToCart(ctx, userID, item)
}

func (s *CartService) GetCart(ctx context.Context, userID primitive.ObjectID) (*model.Cart, error) {
	return s.repo.GetCart(ctx, userID)
}
