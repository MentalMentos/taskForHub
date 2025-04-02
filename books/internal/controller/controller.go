package controller

import (
	"github.com/MentalMentos/taskForHub/books/internal/model"
	"github.com/MentalMentos/taskForHub/books/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BookController struct {
	repo *repository.BookRepository
}

func NewBookController(repo *repository.BookRepository) *BookController {
	return &BookController{repo: repo}
}

func (c *BookController) CreateBook(ctx *gin.Context) {
	var book model.Book
	if err := ctx.ShouldBindJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := c.repo.Create(ctx, &book)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func (c *BookController) GetAllBooks(ctx *gin.Context) {
	books, err := c.repo.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}

	ctx.JSON(http.StatusOK, books)
}
