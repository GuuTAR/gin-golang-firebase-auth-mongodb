package models

type AuthMeResponse struct {
	UID           string `json:"uid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type AuthSignInRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthSignInResponse struct {
	Token     string `json:"token"`      // use as: Authorization: Bearer <token>
	ExpiresIn string `json:"expires_in"` // seconds until expiry (usually 3600)
	UID       string `json:"uid"`
	Email     string `json:"email"`
}
