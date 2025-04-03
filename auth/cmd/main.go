package main

import (
	"github.com/MentalMentos/taskForHub/auth/internal/config"
	"github.com/MentalMentos/taskForHub/auth/internal/controller"
	"github.com/MentalMentos/taskForHub/auth/internal/repository"
	"github.com/MentalMentos/taskForHub/auth/internal/service"
	zaplogger "github.com/MentalMentos/taskForHub/auth/pkg/logger/zap"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Welcome Home!")
	})
	router.GET("/ip", func(c *gin.Context) {
		// Получаем IP клиента
		clientIP := c.ClientIP()
		c.JSON(200, gin.H{"ip": clientIP})
	})

	log := zaplogger.New()
	cfg := config.New(log)
	db := config.DataBaseConnection(cfg)

	authRepository := repository.NewRepository(db, log)
	authService := service.New(authRepository, log)
	authController := controller.NewAuthController(authService, log)

	authRoutes := router.Group("/auth_v1")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/refresh", authController.RefreshToken)
	}
	if err := router.Run("localhost:8081"); err != nil {
		log.Fatal("Main", "Failed to start server")
	}
}
