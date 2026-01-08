package main

import (
	"log/slog"
	"os"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	setupSlog()

	controller := controllers.New()

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
