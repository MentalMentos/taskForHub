package main

import (
	"github.com/MentalMentos/taskForHub/book/internal/config"
	"github.com/MentalMentos/taskForHub/book/internal/controller"
	"github.com/MentalMentos/taskForHub/book/internal/repository"
	"github.com/MentalMentos/taskForHub/book/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
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

	r.Run(":8081")
}
