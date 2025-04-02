package request

import "go.mongodb.org/mongo-driver/bson/primitive"

// AddToCartRequest модель запроса для добавления товара в корзину
type AddToCartRequest struct {
	UserID    primitive.ObjectID `json:"user_id" binding:"required"`
	ProductID primitive.ObjectID `json:"product_id" binding:"required"`
	Quantity  int                `json:"quantity" binding:"required"`
}

type GetCartRequest struct {
	UserID primitive.ObjectID `json:"user_id" binding:"required"`
}
