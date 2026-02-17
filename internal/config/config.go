package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URI     string
	JWT_SECRET string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var config Config
	config.DB_URI = os.Getenv("DATABASE_URL")
	config.JWT_SECRET = os.Getenv("JWT_SECRET")

	return config
}
