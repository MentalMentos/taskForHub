package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
}

type Cart struct {
	Items []CartItem `bson:"items" json:"items"`
}
