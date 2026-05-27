package config

import "os"

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Port       string
	Env        string
	MongoURI   string
	MongoDB    string
}

// Load reads configuration from environment variables, falling back to defaults.
func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		Env:      getEnv("ENV", "development"),
		MongoURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGODB_DB", "app"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
