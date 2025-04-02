package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Book представляет книгу в MongoDB
type Book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Author string             `bson:"author"`
	Title  string             `bson:"title"`
}
