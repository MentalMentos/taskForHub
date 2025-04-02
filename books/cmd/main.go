package main

import (
	"github.com/MentalMentos/taskForHub/books/internal/config"
	"github.com/MentalMentos/taskForHub/books/internal/controller"
	"github.com/MentalMentos/taskForHub/books/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных MongoDB
	db := config.DataBaseConnection()

	// Инициализация репозитория и контроллера
	repo := repository.NewBookRepository(db)
	controller := controller.NewBookController(repo)

	// Инициализация Gin
	r := gin.Default()

	// Роуты
	r.POST("/books", controller.CreateBook)
	r.GET("/books", controller.GetAllBooks)

	// Запуск сервера
	r.Run(":8082")
}
