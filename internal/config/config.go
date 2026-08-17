package config

import (
	"os"
	"strings"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/pkg/auth"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Port     string
	Env      string
	MongoURI string
	MongoDB  string

	// Firebase — individual fields from the service-account JSON.
	// Leave all empty on GCP / Firebase Hosting / Cloud Run to use
	// Application Default Credentials automatically.
	Firebase auth.Credentials

	// FirebaseWebAPIKey is the client-facing Web API key from the Firebase Console
	// (Project Settings → General → Web API key).
	// Only required for the POST /auth/signin sign-in helper endpoint.
	FirebaseWebAPIKey string

	// CORSAllowedOrigins are the origins allowed to make cross-origin requests,
	// e.g. "https://app.example.com,http://localhost:3000". Defaults to "*"
	// (any origin) for ease of local development.
	CORSAllowedOrigins []string
}

// Load reads configuration from environment variables, falling back to defaults.
func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		Env:      getEnv("ENV", "development"),
		MongoURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGODB_DB", "app"),

		Firebase: auth.Credentials{
			ProjectID:    getEnv("FIREBASE_PROJECT_ID", ""),
			PrivateKeyID: getEnv("FIREBASE_PRIVATE_KEY_ID", ""),
			PrivateKey:   getEnv("FIREBASE_PRIVATE_KEY", ""),
			ClientEmail:  getEnv("FIREBASE_CLIENT_EMAIL", ""),
		},
		FirebaseWebAPIKey: getEnv("FIREBASE_WEB_API_KEY", ""),

		CORSAllowedOrigins: getEnvList("CORS_ALLOWED_ORIGINS", []string{"*"}),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// getEnvList reads a comma-separated env var into a trimmed string slice.
func getEnvList(key string, defaultVal []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}

	parts := strings.Split(v, ",")
	list := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			list = append(list, p)
		}
	}
	return list
}
