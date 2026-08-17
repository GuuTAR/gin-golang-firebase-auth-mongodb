package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/middleware"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/models"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/services"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/util"
)

// AuthHandler serves auth-related endpoints.
type AuthHandler struct {
	webAPIKey string
	svc       *services.AuthService
}

func NewAuthHandler(webAPIKey string, svc *services.AuthService) *AuthHandler {
	return &AuthHandler{webAPIKey: webAPIKey, svc: svc}
}

// Me godoc
//
//	@Summary     Get current user
//	@Description Returns the Firebase UID and claims of the authenticated caller.
//	@Tags        auth
//	@Produce     json
//	@Security    BearerAuth
//	@Success     200  {object}  models.Response{body=models.AuthMeResponse}
//	@Failure     401  {object}  models.ResponseWithoutBody
//	@Router      /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	token, ok := middleware.TokenFromContext(c)
	if !ok {
		// Should never reach here if the route is behind FirebaseAuth middleware,
		// but guard defensively.
		util.JSONwithoutBody(c, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// token.Claims contains all JWT payload fields (email, name, picture, …).
	// Expose only the safe subset; add more fields as your app requires.
	emailVerified, _ := token.Claims["email_verified"].(bool)
	util.JSON(c, http.StatusOK, "Get user information successfully", models.AuthMeResponse{
		UID:           token.UID,
		Email:         claimString(token.Claims, "email"),
		EmailVerified: emailVerified,
		Name:          claimString(token.Claims, "name"),
		Picture:       claimString(token.Claims, "picture"),
	})
}

// Token godoc
//
//	@Summary     Sign in and get a Bearer token
//	@Description Exchanges email + password for a Firebase ID token.
//	@Description Use the returned value as: Authorization: Bearer <token>
//	@Tags        auth
//	@Accept      json
//	@Produce     json
//	@Param       body  body      models.AuthSignInRequest  true  "Credentials"
//	@Success     200   {object}  models.Response{body=models.AuthSignInResponse}
//	@Failure     400   {object}  models.ResponseWithoutBody
//	@Failure     401   {object}  models.ResponseWithoutBody
//	@Failure     500   {object}  models.ResponseWithoutBody
//	@Router      /api/v1/auth/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	if h.webAPIKey == "" {
		util.JSONwithoutBody(c, http.StatusInternalServerError,
			"FIREBASE_WEB_API_KEY is not configured (Firebase Console → Project Settings → General → Web API key)")
		return
	}

	var req models.AuthSignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.JSONwithoutBody(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.svc.SignIn(h.webAPIKey, req.Email, req.Password)
	if err != nil {
		util.JSONwithoutBody(c, http.StatusUnauthorized, err.Error())
		return
	}

	util.JSON(c, http.StatusOK, "Signed in successfully", models.AuthSignInResponse{
		Token:     result.IDToken,
		ExpiresIn: result.ExpiresIn,
		UID:       result.LocalID,
		Email:     result.Email,
	})
}

func claimString(claims map[string]interface{}, key string) string {
	v, _ := claims[key].(string)
	return v
}
