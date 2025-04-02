package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
)

type AuthRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

var authServiceURL = "http://localhost:8081"
var productServiceURL = "http://localhost:8082"
var cartServiceURL = "http://localhost:8083"

var jwtSecret = []byte("secret_key")

func main() {
	r := gin.Default()

	// Группа маршрутов для авторизации
	auth := r.Group("/auth_v1")
	{
		auth.POST("/register", registerHandler)
		auth.POST("/login", loginHandler)
	}

	// Группа маршрутов для сервисов
	authorized := r.Group("/")
	authorized.Use(authMiddleware())
	{
		authorized.GET("/products", proxyHandler(productServiceURL))
		authorized.GET("/cart", proxyHandler(cartServiceURL))
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Main", "Failed to start server")
	}
}

func registerHandler(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	client := resty.New()
	var authResp AuthResponse

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&authResp).
		Post(authServiceURL + "/auth_v1/register")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to auth service"})
		return
	}

	if resp.IsError() {
		c.JSON(resp.StatusCode(), gin.H{"error": "Registration failed"})
		return
	}

	// Возвращаем полученный access_token клиенту
	c.JSON(http.StatusOK, authResp)
}

func loginHandler(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	client := resty.New()
	var authResp AuthResponse

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&authResp).
		Post(authServiceURL + "/auth_v1/login")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to auth service"})
		return
	}

	if resp.IsError() {
		c.JSON(resp.StatusCode(), gin.H{"error": "Invalid credentials"})
		return
	}

	// Возвращаем токен клиенту
	c.JSON(http.StatusOK, authResp)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func proxyHandler(serviceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", c.GetHeader("Authorization")).
			Get(serviceURL + c.Request.URL.Path)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to service"})
			return
		}

		c.Data(resp.StatusCode(), "application/json", resp.Body())
	}
}
