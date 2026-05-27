package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/middleware"
)

// AuthHandler serves auth-related endpoints.
type AuthHandler struct{}

func NewAuthHandler() *AuthHandler { return &AuthHandler{} }

// Me godoc
//
//	@Summary     Get current user
//	@Description Returns the Firebase UID and claims of the authenticated caller.
//	@Tags        auth
//	@Produce     json
//	@Security    BearerAuth
//	@Success     200  {object}  map[string]any
//	@Failure     401  {object}  map[string]string
//	@Router      /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	token, ok := middleware.TokenFromContext(c)
	if !ok {
		// Should never reach here if the route is behind FirebaseAuth middleware,
		// but guard defensively.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	// token.Claims contains all JWT payload fields (email, name, picture, …).
	// Expose only the safe subset; add more fields as your app requires.
	c.JSON(http.StatusOK, gin.H{
		"uid":            token.UID,
		"email":          token.Claims["email"],
		"email_verified": token.Claims["email_verified"],
		"name":           token.Claims["name"],
		"picture":        token.Claims["picture"],
	})
}
