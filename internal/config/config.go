package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
	// Add other config vars as needed, e.g., for notifications (SMTP settings, etc.)
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // Ignore error if .env not found

	return &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"), // Default to ":8080" if empty
	}, nil
}
