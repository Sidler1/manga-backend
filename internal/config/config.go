package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
	JWTSecret     string // Added for JWT
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // Ignore error if .env not found

	return &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		JWTSecret:     os.Getenv("JWT_SECRET"), // Set in .env
	}, nil
}
