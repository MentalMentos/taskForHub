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
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   struct {
		AccessToken string `json:"access_token"`
		UserID      string `json:"id"`
	} `json:"data"`
}

type AddToCartRequest struct {
	ItemID   string `json:"item_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

type GetCartRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type CartResponse struct {
	UserID string `json:"user_id"`
	ItemID string `json:"item_id"`
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

	authorized := r.Group("")
	authorized.Use(authMiddleware())
	{
		// Продукты
		authorized.GET("/books", proxyGET(productServiceURL))
		authorized.POST("/books", proxyPOST(productServiceURL))

		// Исправленные пути для корзины
		authorized.GET("/cart", getCartHandler)    // для получения корзины
		authorized.POST("/cart", addToCartHandler) // для добавления в корзину
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server")
	}
}

// Обработчик регистрации
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

func getCartHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id"})
		return
	}

	req := GetCartRequest{
		UserID: userID.(string),
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", c.GetHeader("Authorization")).
		SetBody(req).
		Get(cartServiceURL + "/cart-get")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to cart service"})
		return
	}

	c.Data(resp.StatusCode(), "application/json", resp.Body())
}

func addToCartHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id"})
		return
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	reqData := map[string]interface{}{
		"user_id":  userID.(string),
		"item_id":  req.ItemID,
		"quantity": req.Quantity,
	}
	
	client := resty.New()
	var cartresp CartResponse
	_, err := client.R().
		SetHeader("Authorization", c.GetHeader("Authorization")).
		SetHeader("Content-Type", "application/json").
		SetBody(reqData).
		SetResult(&cartresp).
		Post(cartServiceURL + "/cart-add")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to cart service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_cart": cartresp.UserID,
		"item_id":   cartresp.ItemID,
	})
}
