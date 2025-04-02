package model

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
}

type Cart struct {
	mgm.DefaultModel `bson:",inline"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	Items            []CartItem         `json:"items" bson:"items"`
}
