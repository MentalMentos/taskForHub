package controller

import (
	"github.com/MentalMentos/taskForHub/books/internal/model"
	"github.com/MentalMentos/taskForHub/books/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	service *service.BookService
}

func NewBookController(service *service.BookService) *BookController {
	return &BookController{service: service}
}

func (c *BookController) CreateBook(ctx *gin.Context) {
	var book model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.CreateBook(ctx, &book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Book added successfully")
}

func (c *BookController) GetAllBooks(ctx *gin.Context) {
	books, err := c.service.GetAllBooks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}
