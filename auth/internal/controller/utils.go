package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"task0325/auth/internal/data/response"
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

func GetClientIP(c *gin.Context) string {
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		clientIP := strings.TrimSpace(ips[0])
		fmt.Printf("X-Forwarded-For: %s, Resolved IP: %s\n", xForwardedFor, clientIP)
		return clientIP
	}

	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		fmt.Printf("X-Real-IP: %s\n", xRealIP)
		return xRealIP
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		fmt.Printf("RemoteAddr (raw): %s, Error: %v\n", c.Request.RemoteAddr, err)
		return c.Request.RemoteAddr
	}
	fmt.Printf("RemoteAddr (parsed): %s\n", ip)
	return ip
}
