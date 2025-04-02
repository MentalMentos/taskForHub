package main

import (
	"github.com/MentalMentos/taskForHub/books/internal/config"
	"github.com/MentalMentos/taskForHub/books/internal/controller"
	"github.com/MentalMentos/taskForHub/books/internal/repository"
	"github.com/MentalMentos/taskForHub/books/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/kamva/mgm/v3"
)

func main() {
	// Инициализируем MongoDB
	config.InitMongoDB()

	r := gin.Default()

	bookRepo := repository.NewBookRepository()
	bookService := service.NewBookService(bookRepo)
	bookController := controller.NewBookController(bookService)

	r.POST("/books", bookController.CreateBook)
	r.GET("/books", bookController.GetAllBooks)

	r.Run("localhost:8082")
}
