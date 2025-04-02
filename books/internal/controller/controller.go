package controller

import (
	"net/http"

	"github.com/MentalMentos/taskForHub/book/internal/model"
	"github.com/MentalMentos/taskForHub/book/internal/service"
	"github.com/gin-gonic/gin"
)

type BookController struct {
	service *service.BookService
}

func NewBookController(service *service.BookService) *BookController {
	return &BookController{service: service}
}

// @Summary Добавить книгу
// @Tags books
// @Accept json
// @Produce json
// @Param book body model.Book true "Book"
// @Success 200 {string} string "Book added successfully"
// @Router /books [post]
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

// @Summary Получить все книги
// @Tags books
// @Produce json
// @Success 200 {array} model.Book
// @Router /books [get]
func (c *BookController) GetAllBooks(ctx *gin.Context) {
	books, err := c.service.GetAllBooks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}
