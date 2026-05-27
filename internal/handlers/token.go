package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const firebaseSignInURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s"

// TokenHandler issues Firebase ID tokens for development / testing.
// In production, the client app should sign in directly with the Firebase SDK.
type TokenHandler struct {
	webAPIKey string
}

func NewTokenHandler(webAPIKey string) *TokenHandler {
	return &TokenHandler{webAPIKey: webAPIKey}
}

// ── HTTP types ────────────────────────────────────────────────────────────────

type signInRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type signInResponse struct {
	Token     string `json:"token"`      // use as: Authorization: Bearer <token>
	ExpiresIn string `json:"expires_in"` // seconds until expiry (usually 3600)
	UID       string `json:"uid"`
	Email     string `json:"email"`
}

// ── Handler ───────────────────────────────────────────────────────────────────

// Token godoc
//
//	@Summary     Sign in and get a Bearer token
//	@Description Exchanges email + password for a Firebase ID token.
//	@Description Use the returned value as: Authorization: Bearer <token>
//	@Tags        auth
//	@Accept      json
//	@Produce     json
//	@Param       body  body      signInRequest  true  "Credentials"
//	@Success     200   {object}  signInResponse
//	@Failure     400   {object}  map[string]string
//	@Failure     401   {object}  map[string]string
//	@Failure     500   {object}  map[string]string
//	@Router      /api/v1/auth/token [post]
func (h *TokenHandler) Token(c *gin.Context) {
	if h.webAPIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "FIREBASE_WEB_API_KEY is not configured (Firebase Console → Project Settings → General → Web API key)",
		})
		return
	}

	var req signInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.signIn(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token:     result.IDToken,
		ExpiresIn: result.ExpiresIn,
		UID:       result.LocalID,
		Email:     result.Email,
	})
}

// ── Firebase Identity Toolkit ─────────────────────────────────────────────────

type firebaseSignInPayload struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type firebaseSignInResult struct {
	IDToken   string         `json:"idToken"`
	LocalID   string         `json:"localId"`
	Email     string         `json:"email"`
	ExpiresIn string         `json:"expiresIn"`
	Error     *firebaseError `json:"error,omitempty"`
}

type firebaseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ── Internal helpers ──────────────────────────────────────────────────────────

// signIn calls the Firebase REST API and returns the token result, or an error
// if Firebase rejects the credentials.
func (h *TokenHandler) signIn(ctx context.Context, email, password string) (*firebaseSignInResult, error) {
	result, err := h.callFirebase(ctx, email, password)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("%s", result.Error.Message)
	}
	return result, nil
}

// callFirebase performs the raw HTTP call to the Identity Toolkit sign-in endpoint.
func (h *TokenHandler) callFirebase(ctx context.Context, email, password string) (*firebaseSignInResult, error) {
	payload, err := json.Marshal(firebaseSignInPayload{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	url := fmt.Sprintf(firebaseSignInURL, h.webAPIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("firebase request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var result firebaseSignInResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &result, nil
}
