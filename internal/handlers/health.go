package handlers

import (
	"net/http"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/util"
	"github.com/gin-gonic/gin"
)

// HealthHandler serves the /health endpoint (no auth required).
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

// Check godoc
//
//	@Summary     Health check
//	@Description Returns ok when the server is running.
//	@Tags        health
//	@Produce     json
//	@Success     200  {object}  models.ResponseWithoutBody
//	@Router      /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	util.JSONwithoutBody(c, http.StatusOK, "ok")
}
