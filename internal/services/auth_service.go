// Package services contains business logic, independent of HTTP concerns.
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/models"
)

const firebaseSignInURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s"

// AuthService implements the business logic for managing authentication.
type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// signIn calls the Firebase REST API and returns the token result, or an error
// if Firebase rejects the credentials.
func (s *AuthService) SignIn(webAPIKey, email, password string) (*models.FirebaseSignInResult, error) {
	result, err := callFirebase(webAPIKey, email, password)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("%s", result.Error.Message)
	}
	return result, nil
}

// callFirebase performs the raw HTTP call to the Identity Toolkit sign-in endpoint.
func callFirebase(webAPIKey, email, password string) (*models.FirebaseSignInResult, error) {
	payload, err := json.Marshal(models.FirebaseSignInPayload{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	url := fmt.Sprintf(firebaseSignInURL, webAPIKey)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
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

	var result models.FirebaseSignInResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &result, nil
}
