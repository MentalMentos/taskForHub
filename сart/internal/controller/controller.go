package controller

import (
	"github.com/MentalMentos/taskForHub/сart/internal/model"
	"github.com/MentalMentos/taskForHub/сart/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartController struct {
	service *service.CartService
}

func NewCartController(service *service.CartService) *CartController {
	return &CartController{service: service}
}

func (c *CartController) GetCart(ctx *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	cart, err := c.service.GetCart(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var item model.CartItem
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := c.service.AddToCart(ctx, userID, item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "item added"})
}

func (c *CartController) RemoveFromCart(ctx *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	productID, err := primitive.ObjectIDFromHex(ctx.Param("product_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := c.service.RemoveFromCart(ctx, userID, productID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
