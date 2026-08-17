// Package middleware contains Gin middleware used across the application.
package middleware

import (
	"net/http"
	"strings"

	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/util"
	"github.com/gin-gonic/gin"
)

// Context keys for values set by FirebaseAuth.
const (
	ContextKeyUID   = "firebase_uid"
	ContextKeyToken = "firebase_token"
)

// FirebaseAuth returns a Gin middleware that validates Firebase ID tokens sent
// in the Authorization header as "Bearer <token>".
//
// On success it stores two values in the Gin context:
//   - ContextKeyUID   (string)              — the Firebase UID
//   - ContextKeyToken (*firebaseAuth.Token) — the full verified token
//
// Use TokenFromContext / UIDFromContext to retrieve them in handlers.
func FirebaseAuth(client *firebaseAuth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			util.JSONwithoutBody(c, http.StatusUnauthorized, "missing or malformed Authorization header")
			c.Abort()
			return
		}

		idToken := strings.TrimPrefix(header, "Bearer ")

		token, err := client.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			util.JSONwithoutBody(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextKeyUID, token.UID)
		c.Set(ContextKeyToken, token)
		c.Next()
	}
}

// TokenFromContext retrieves the verified Firebase token stored by FirebaseAuth.
// Returns nil, false if the middleware has not run for this request.
func TokenFromContext(c *gin.Context) (*firebaseAuth.Token, bool) {
	val, exists := c.Get(ContextKeyToken)
	if !exists {
		return nil, false
	}
	token, ok := val.(*firebaseAuth.Token)
	return token, ok
}

// UIDFromContext retrieves the Firebase UID stored by FirebaseAuth.
// Returns "", false if the middleware has not run for this request.
func UIDFromContext(c *gin.Context) (string, bool) {
	val, exists := c.Get(ContextKeyUID)
	if !exists {
		return "", false
	}
	uid, ok := val.(string)
	return uid, ok
}
