package main

import (
	"fmt"
	_ "github.com/MentalMentos/taskForHub/api-gateway/docs"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

type AuthRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		AccessToken string `json:"access_token"`
		UserID      string `json:"id"`
	} `json:"data"`
}

var authServiceURL = "http://localhost:8081"
var productServiceURL = "http://localhost:8082"

var jwtSecret = []byte("secret_key")

// @title API Gateway
// @version 1.0
// @description API Gateway для авторизации и управления книгами.
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth_v1")
	{
		auth.POST("/register", registerHandler)
		auth.POST("/login", loginHandler)
	}

	authorized := r.Group("")
	authorized.Use(authMiddleware())
	{
		// Продукты
		authorized.GET("/books", proxyGET(productServiceURL))
		authorized.POST("/books", proxyPOST(productServiceURL))
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server")
	}
}

// registerHandler регистрирует нового пользователя
// @Summary Регистрация
// @Description Регистрирует нового пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Данные пользователя"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /auth_v1/register [post]
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

	if authResp.Data.AccessToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Empty access token from auth service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": authResp.Data.AccessToken, "user_id": authResp.Data.UserID})
}

// loginHandler авторизует пользователя
// @Summary Авторизация
// @Description Авторизует пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Данные пользователя"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /auth_v1/login [post]
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
		c.JSON(resp.StatusCode(), gin.H{"error": "Registration failed"})
		return
	}

	fmt.Printf("Auth Service Response: %+v\n", authResp)

	if authResp.Data.UserID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Empty user ID from auth service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": authResp.Data.AccessToken,
		"user_id":      authResp.Data.UserID,
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		fmt.Println("Parsed Token:", tokenString)

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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id in token"})
			c.Abort()
			return
		}

		fmt.Println("User ID set in context:", userID)

		c.Set("user_id", userID)
		c.Next()
	}
}

// proxyGET проксирует GET-запросы к сервису книг
// @Summary Получение списка книг
// @Description Получает список книг из сервиса товаров
// @Tags Books
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string "Failed to connect to service"
// @Router /books [get]
func proxyGET(serviceURL string) gin.HandlerFunc {
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

// proxyPOST проксирует POST-запросы к сервису книг
// @Summary Добавление книги
// @Description Добавляет книгу в сервис товаров
// @Tags Books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Данные книги"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string "Failed to connect to service"
// @Router /books [post]
func proxyPOST(serviceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id"})
			return
		}

		fmt.Println("User ID:", userID.(string))

		client := resty.New()
		resp, err := client.R().
			SetHeader("Authorization", c.GetHeader("Authorization")).
			SetHeader("user_id", userID.(string)).
			SetHeader("Content-Type", "application/json").
			SetBody(c.Request.Body).
			Post(serviceURL + c.Request.URL.Path)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to service"})
			return
		}

		c.Data(resp.StatusCode(), "application/json", resp.Body())
	}
}
