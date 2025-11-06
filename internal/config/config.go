package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL         string
	Port                string
	JWTSecret           string
	BearerTokenDuration string
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		Port:                getEnv("PORT", ""),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		BearerTokenDuration: getEnv("BEARER_TOKEN_DURATION", "168h"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.BearerTokenDuration == "" {
		return fmt.Errorf("BEARER_TOKEN_DURATION is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
