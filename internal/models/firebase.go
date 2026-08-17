package models

// ── Firebase Identity Toolkit ─────────────────────────────────────────────────

type FirebaseSignInPayload struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type FirebaseSignInResult struct {
	IDToken   string         `json:"idToken"`
	LocalID   string         `json:"localId"`
	Email     string         `json:"email"`
	ExpiresIn string         `json:"expiresIn"`
	Error     *FirebaseError `json:"error,omitempty"`
}

type FirebaseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
