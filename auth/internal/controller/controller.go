package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task0325/auth/internal/data/request"
	"task0325/auth/internal/service"
	"task0325/auth/pkg/logger"
)

type AuthController struct {
	authService service.Service
	logger      logger.Logger
}

func NewAuthController(authService *service.Service, logger logger.Logger) *AuthController {
	return &AuthController{
		authService: *authService,
		logger:      logger,
	}
}

func (controller *AuthController) Register(c *gin.Context) {
	var userRequest request.RegisterUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		HandleError(c, &ApiError{Code: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	userRequest.IP = GetClientIP(c)
	authResp, err := controller.authService.Register(c, userRequest)
	if err != nil {
		HandleError(c, err)
		return
	}

	JsonResponse(c, http.StatusOK, "Registration successful", authResp)
}

func (controller *AuthController) Login(c *gin.Context) {
	var userRequest request.LoginRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		HandleError(c, &ApiError{Code: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	userRequest.IP = GetClientIP(c)
	authResp, err := controller.authService.Login(c, userRequest)
	if err != nil {
		HandleError(c, err)
		return
	}

	JsonResponse(c, http.StatusOK, "Auth successful", authResp)
}

func (controller *AuthController) UpdatePassword(c *gin.Context) {
	var userRequest request.UpdateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		HandleError(c, &ApiError{Code: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	userRequest.IP = GetClientIP(c)
	authResp, err := controller.authService.UpdatePassword(c, userRequest)
	if err != nil {
		HandleError(c, err)
		return
	}

	JsonResponse(c, http.StatusOK, "Password updated successful", authResp)
}

func (controller *AuthController) RefreshToken(c *gin.Context) {
	var userRequest request.UpdateTokenRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		HandleError(c, &ApiError{Code: http.StatusBadRequest, Message: "Invalid request payload"})
		return
	}

	authResp, err := controller.authService.GetAccessToken(c, userRequest.RefreshToken)
	if err != nil {
		HandleError(c, err)
		return
	}

	JsonResponse(c, http.StatusOK, "Token refreshed successful", authResp)
}
