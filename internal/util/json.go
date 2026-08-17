package util

import (
	"github.com/gin-gonic/gin"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/models"
)

// JSON writes status with message and body.
func JSON(c *gin.Context, status int, message string, body any) {
	c.JSON(status, models.Response{Message: message, Body: body})
}

// JSON writes status with message and no body.
func JSONwithoutBody(c *gin.Context, status int, message string) {
	c.JSON(status, models.ResponseWithoutBody{Message: message})
}
