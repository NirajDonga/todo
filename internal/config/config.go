package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URI     string
	JWT_SECRET string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("loading .env: %w", err)
	}

	cfg := Config{
		DB_URI:     os.Getenv("DATABASE_URL"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
	}

	if cfg.DB_URI == "" {
		return Config{}, fmt.Errorf("environment variable DATABASE_URL is required")
	}

	return cfg, nil
}
