package config

import (
	"log"
	"sync"

	"github.com/kamva/mgm/v3"
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
