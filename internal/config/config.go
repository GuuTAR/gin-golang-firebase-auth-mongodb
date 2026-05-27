package config

import (
	"os"

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
	// Only required for the POST /auth/token sign-in helper endpoint.
	FirebaseWebAPIKey string
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
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
