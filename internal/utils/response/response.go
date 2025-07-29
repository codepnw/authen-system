package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, gin.H{"message": message, "data": data})
}

func BadRequest(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"message": message, "error": err.Error()})
}

func Unauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized", "error": err.Error()})
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
