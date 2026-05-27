package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler serves the /health endpoint (no auth required).
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
