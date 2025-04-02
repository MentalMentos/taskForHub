package main

import (
	"github.com/MentalMentos/taskForHub/сart/internal/config"
	"github.com/MentalMentos/taskForHub/сart/internal/controller"
	"github.com/MentalMentos/taskForHub/сart/internal/repository"
	"github.com/MentalMentos/taskForHub/сart/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Подключение к MongoDB
	err := mgm.SetDefaultConfig(nil, "cart_db", options.Client().ApplyURI("mongodb://user:password@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	config.InitMongoDB()

	repo := repository.NewCartRepository()
	svc := service.NewCartService(repo)
	ctrl := controller.NewCartController(svc)

	api := r.Group("/cart_v1")
	{
		api.GET("/:user_id", ctrl.GetCart)
		api.POST("/:user_id/add", ctrl.AddToCart)
		api.DELETE("/:user_id/remove/:product_id", ctrl.RemoveFromCart)
	}

	r.Run("localhost:8083")
}
