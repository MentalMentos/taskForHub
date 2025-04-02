package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Определите константы для конфигурации MongoDB
const (
	MONGO_HOST     = "localhost"
	MONGO_PORT     = "27017"
	MONGO_USER     = "user"
	MONGO_PASSWORD = "1234"
	MONGO_DBNAME   = "mongo"
)

// DataBaseConnection создает подключение к базе данных MongoDB
func DataBaseConnection() *mongo.Database {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin&authMechanism=SCRAM-SHA-256",
		MONGO_USER, MONGO_PASSWORD, MONGO_HOST, MONGO_PORT, MONGO_DBNAME)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Ошибка при подключении к MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ошибка при пинге MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client.Database(MONGO_DBNAME)
}
