package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/controllers"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/services"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
		os.Exit(1)
	}

	setupSlog()

	repos, err := repositories.New(os.Getenv("DB_DSN"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	jwtService := jwt.New(os.Getenv("JWT_KEY"))
	services := services.New(repos)

	controller := controllers.New(services, jwtService)

	r := gin.Default()

	controller.Routes(r)

	if err := r.Run(); err != nil {
		slog.Error(err.Error())
	}
}

// setupSlog setup the slog library provided by Go standard library for
// structured logging
func setupSlog() {
	loggerOpts := &slog.HandlerOptions{}
	if os.Getenv("ENVIRONMENT") == "DEVELOPMENT" || os.Getenv("ENVIRONMENT") == "TESTING" {
		loggerOpts.Level = slog.LevelDebug
	} else {
		loggerOpts.Level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOpts))
	slog.SetDefault(logger)
}
