package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Author string             `bson:"author"`
	Title  string             `bson:"title"`
}
