package config

import (
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() {
	err := mgm.SetDefaultConfig(nil, "book_db", options.Client().ApplyURI("mongodb://user:password@mongo:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
}
