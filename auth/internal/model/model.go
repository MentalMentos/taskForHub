package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"` // MongoDB ID
	Name     string             `json:"username" binding:"required" bson:"username"`
	Email    string             `json:"email" binding:"required" bson:"email"`
	Password string             `json:"password" binding:"required" bson:"password"`
}

type UserApi struct {
	ID    string `json:"id"`
	Name  string `json:"username" binding:"required" bson:"username"`
	Email string `json:"email" binding:"required" bson:"email"`
}
