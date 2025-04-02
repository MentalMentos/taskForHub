package controller

import (
	request "github.com/MentalMentos/taskForHub/cart/data"
	"github.com/MentalMentos/taskForHub/cart/internal/model"
	"github.com/MentalMentos/taskForHub/cart/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CartController - структура контроллера корзины
type CartController struct {
	service service.CartService
}

// NewCartController - конструктор для CartController
func NewCartController(service *service.CartService) *CartController {
	return &CartController{service: *service}
}

// AddToCartHandler - обработчик добавления товара в корзину
func (c *CartController) AddToCartHandler(ctx *gin.Context) {
	var req request.AddToCartRequest
	// Читаем тело запроса и связываем с моделью
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Добавляем товар в корзину
	err := c.service.AddToCart(ctx, req.UserID, model.CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

func (c *CartController) GetCartHandler(ctx *gin.Context) {
	var req request.GetCartRequest
	// Получаем корзину пользователя
	cart, err := c.service.GetCart(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}
