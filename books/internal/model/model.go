package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Book структура книги
type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Author      string             `json:"author" binding:"required" bson:"author"`
	Title       string             `json:"title" binding:"required" bson:"title"`
	Price       float64            `json:"price" bson:"price"`
	Description string             `json:"description" bson:"description"`
}
