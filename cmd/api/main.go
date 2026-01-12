package main

import (
	"log/slog"
	"os"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/config"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/controllers"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/services"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/jwt"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           FAQ API
// @version         1.0
// @description     FAQ Management System

// @contact.name   Ahmad Abdelrazik
// @contact.email  ahmad.abdelrazik.swe@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	cfg := config.Load()

	setupSlog(cfg.Environment)

	repos, err := repositories.New(cfg.DSN)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	jwtService := jwt.New(cfg.JwtKey)
	services := services.New(repos)

	controller := controllers.New(services, jwtService)

	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	controller.Routes(r)

	if err := r.Run(cfg.Port); err != nil {
		slog.Error(err.Error())
	}
}

// setupSlog setup the slog library provided by Go standard library for
// structured logging
func setupSlog(environment string) {
	loggerOpts := &slog.HandlerOptions{}
	if environment == "DEVELOPMENT" || environment == "TESTING" {
		loggerOpts.Level = slog.LevelDebug
	} else {
		loggerOpts.Level = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOpts))
	slog.SetDefault(logger)
}
