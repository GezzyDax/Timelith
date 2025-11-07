package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	DatabaseURL          string
	RedisURL             string
	TelegramAppID        string
	TelegramAppHash      string
	APIKey               string
	SessionEncryptionKey string
	Environment          string
}

func Load() (*Config, error) {
	// Load .env file if it exists (for development)
	_ = godotenv.Load()

	cfg := &Config{
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		RedisURL:             getEnv("REDIS_URL", "redis://localhost:6379/0"),
		TelegramAppID:        getEnv("TELEGRAM_APP_ID", ""),
		TelegramAppHash:      getEnv("TELEGRAM_APP_HASH", ""),
		APIKey:               getEnv("GO_API_KEY", ""),
		SessionEncryptionKey: getEnv("SESSION_ENCRYPTION_KEY", ""),
		Environment:          getEnv("ENVIRONMENT", "production"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.TelegramAppID == "" {
		return fmt.Errorf("TELEGRAM_APP_ID is required")
	}
	if c.TelegramAppHash == "" {
		return fmt.Errorf("TELEGRAM_APP_HASH is required")
	}
	if c.APIKey == "" {
		return fmt.Errorf("GO_API_KEY is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
