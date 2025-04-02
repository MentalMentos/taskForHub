package config

import (
	"context"
	"fmt"
	"github.com/MentalMentos/taskForHub/auth/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

// Определите константы для конфигурации MongoDB
const (
	MONGO_HOST     = "localhost"
	MONGO_PORT     = "27017"
	MONGO_USER     = "user"
	MONGO_PASSWORD = "1234"
	MONGO_DBNAME   = "mongo"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

var (
	config Config
	once   sync.Once
)

func New(logger logger.Logger) *Config {
	once.Do(func() {
		config = Config{
			Host:     MONGO_HOST,
			Port:     MONGO_PORT,
			Username: MONGO_USER,
			Password: MONGO_PASSWORD,
			DBName:   MONGO_DBNAME,
		}
	})
	logger.Info("Config", "Config init")
	return &config
}

func DataBaseConnection(cfg *Config) *mongo.Database {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin&authMechanism=SCRAM-SHA-256",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

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
	return client.Database(cfg.DBName)
}
