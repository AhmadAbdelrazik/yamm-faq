package config

import (
	"log"
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
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
		os.Exit(1)
	}

	return &Config{
		DSN:         os.Getenv("DB_DSN"),
		Port:        os.Getenv("PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
		JwtKey:      os.Getenv("JWT_KEY"),
	}
}
