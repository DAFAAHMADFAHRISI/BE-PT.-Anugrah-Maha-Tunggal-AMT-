package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, httpCode int, message string, data interface{}) {
	c.JSON(httpCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, httpCode int, message string, err interface{}) {
	c.JSON(httpCode, APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func BadRequest(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusBadRequest, message, err)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

func InternalServerError(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusInternalServerError, message, err)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}
