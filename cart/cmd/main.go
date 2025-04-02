package main

import (
	"github.com/MentalMentos/taskForHub/cart/internal/config"
	"github.com/MentalMentos/taskForHub/cart/internal/controller"
	"github.com/MentalMentos/taskForHub/cart/internal/repository"
	"github.com/MentalMentos/taskForHub/cart/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Инициализация подключения к базе данных MongoDB
	db := config.DataBaseConnection()

	// Инициализация репозитория
	repo := repository.NewCartRepository(db)

	// Инициализация сервиса
	cartService := service.NewCartService(repo)

	// Инициализация контроллера
	cartController := controller.NewCartController(cartService)

	// Инициализация Gin
	r := gin.Default()

	// Роуты для работы с корзиной
	router := gin.Default()

	cartRoutes := router.Group("/cart")
	{
		cartRoutes.GET("/cart-get", cartController.GetCartHandler)
		cartRoutes.POST("/cart-add", cartController.AddToCartHandler)
	}

	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
