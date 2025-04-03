package main

import (
	"github.com/MentalMentos/taskForHub/books/internal/config"
	"github.com/MentalMentos/taskForHub/books/internal/controller"
	"github.com/MentalMentos/taskForHub/books/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.DataBaseConnection()

	repo := repository.NewBookRepository(db)
	controller := controller.NewBookController(repo)

	r := gin.Default()

	r.POST("/books", controller.CreateBook)
	r.GET("/books", controller.GetAllBooks)

	r.Run(":8082")
}
