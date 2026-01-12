package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN         string
	Port        string
	Environment string
	JwtKey      string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DSN:         os.Getenv("DB_DSN"),
		Port:        os.Getenv("PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
		JwtKey:      os.Getenv("JWT_KEY"),
	}
}
