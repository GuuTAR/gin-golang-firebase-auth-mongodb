// Package auth provides a thin wrapper around the Firebase Admin SDK.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Credentials holds individual service-account fields loaded from env vars.
// Each field maps 1-to-1 with a key in the service-account JSON downloaded
// from the Firebase Console.
type Credentials struct {
	ProjectID    string // FIREBASE_PROJECT_ID
	PrivateKeyID string // FIREBASE_PRIVATE_KEY_ID
	PrivateKey   string // FIREBASE_PRIVATE_KEY   (literal \n are unescaped automatically)
	ClientEmail  string // FIREBASE_CLIENT_EMAIL
	ClientID     string // FIREBASE_CLIENT_ID
	CertURL      string // FIREBASE_CLIENT_CERT_URL
}

// IsSet reports whether all required credential fields are populated.
func (c Credentials) IsSet() bool {
	return c.ProjectID != "" &&
		c.PrivateKeyID != "" &&
		c.PrivateKey != "" &&
		c.ClientEmail != ""
}

// serviceAccountJSON is the wire shape expected by google.golang.org/api/option.
type serviceAccountJSON struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

// buildJSON assembles the individual credential fields back into the JSON that
// the Google SDK expects.  Literal \n sequences in PrivateKey are converted to
// real newlines so the value can be stored as a plain env var.
func buildJSON(c Credentials) ([]byte, error) {
	sa := serviceAccountJSON{
		Type:                    "service_account",
		ProjectID:               c.ProjectID,
		PrivateKeyID:            c.PrivateKeyID,
		PrivateKey:              strings.ReplaceAll(c.PrivateKey, `\n`, "\n"),
		ClientEmail:             c.ClientEmail,
		ClientID:                c.ClientID,
		AuthURI:                 "https://accounts.google.com/o/oauth2/auth",
		TokenURI:                "https://oauth2.googleapis.com/token",
		AuthProviderX509CertURL: "https://www.googleapis.com/oauth2/v1/certs",
		ClientX509CertURL:       c.CertURL,
	}
	return json.Marshal(sa)
}

// FirebaseClient wraps the Firebase Auth client.
type FirebaseClient struct {
	Auth *firebaseAuth.Client
}

// NewFirebaseClient initialises the Firebase Admin SDK and returns an Auth client.
//
// Credential resolution order:
//  1. creds — individual service-account fields from FIREBASE_* env vars.
//  2. GOOGLE_APPLICATION_CREDENTIALS env var (path to a service-account file).
//  3. Application Default Credentials — automatic on GCP / Firebase Hosting / Cloud Run.
func NewFirebaseClient(ctx context.Context, creds Credentials) (*FirebaseClient, error) {
	var opts []option.ClientOption

	if creds.IsSet() {
		jsonBytes, err := buildJSON(creds)
		if err != nil {
			return nil, fmt.Errorf("firebase: build credentials JSON: %w", err)
		}
		opts = append(opts, option.WithAuthCredentialsJSON(option.ServiceAccount, jsonBytes))
	}

	var cfg *firebase.Config
	if creds.ProjectID != "" {
		cfg = &firebase.Config{ProjectID: creds.ProjectID}
	}

	app, err := firebase.NewApp(ctx, cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("firebase: init app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase: auth client: %w", err)
	}

	return &FirebaseClient{Auth: authClient}, nil
}
