package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port            string
	MongoURI        string
	DBName          string
	JWTSecretKey    string
	TokenDuration   time.Duration
	RateLimit       int
	RateLimitWindow time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config { // Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found or error loading it, using environment variables or defaults")
	} else {
		log.Println("📄 Successfully loaded .env file")
	}

	config := &Config{
		Port:            getEnv("PORT", ":50051"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:          getEnv("DB_NAME", "usermanagement"),
		JWTSecretKey:    getEnv("JWT_SECRET_KEY", "your-secret-key"),
		TokenDuration:   getDurationEnv("TOKEN_DURATION", 24*time.Hour),
		RateLimit:       getIntEnv("RATE_LIMIT", 5),
		RateLimitWindow: getDurationEnv("RATE_LIMIT_WINDOW", time.Minute),
	}

	return config
}

// Helper functions
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}
