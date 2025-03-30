package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"sync"
	"task0325/auth/pkg/helpers"
	"task0325/auth/pkg/logger"
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

const (
	host     = "localhost"
	user     = "user"
	password = "1234"
	dbName   = "postgres"
)

func New(logger logger.Logger) *Config {
	once.Do(func() {
		config = Config{
			Host:     os.Getenv("PG_HOST"),
			Port:     os.Getenv("PG_PORT"),
			Username: os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			DBName:   os.Getenv("PG_DATABASE_NAME"),
		}
	})
	logger.Info("Config", "Config init")
	return &config
}

func DatabaseConnection(logger logger.Logger) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(helpers.PgPrefix, helpers.PgConnectFailed)
	}
	logger.Info(helpers.PgPrefix, "Database connection done")
	return db
}
