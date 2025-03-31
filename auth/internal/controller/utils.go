package controller

import (
	"github.com/MentalMentos/taskForHub/auth/internal/data/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiError struct {
	Code    int
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

func HandleError(c *gin.Context, err error) {
	if apiErr, ok := err.(*ApiError); ok {
		JsonResponse(c, apiErr.Code, apiErr.Message, nil)
	} else {
		JsonResponse(c, http.StatusInternalServerError, "Internal Server Error", nil)
	}
}

func JsonResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, response.Response{
		Code:   status,
		Status: message,
		Data:   data,
	})
}
