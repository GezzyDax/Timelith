package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string
	APIKey     string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// Telegram
	TelegramAppID   int
	TelegramAppHash string

	// Security
	JWTSecret     string
	EncryptionKey string

	// Environment
	Environment string
}

func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		APIKey:          getEnv("API_KEY", ""),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		TelegramAppHash: getEnv("TELEGRAM_APP_HASH", ""),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		EncryptionKey:   getEnv("ENCRYPTION_KEY", ""),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}

	// Parse TelegramAppID
	appIDStr := getEnv("TELEGRAM_APP_ID", "0")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TELEGRAM_APP_ID: %w", err)
	}
	cfg.TelegramAppID = appID

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.TelegramAppID == 0 || cfg.TelegramAppHash == "" {
		return nil, fmt.Errorf("TELEGRAM_APP_ID and TELEGRAM_APP_HASH are required")
	}
	if cfg.EncryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
