package controllers

import (
	"github.com/AhmadAbdelrazik/yamm_faq/internal/services"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/jwt"
)

// Controller handle the first layer of the application including routes,
// middlewares, and handlers.
type Controller struct {
	Services *services.Services
	jwt      *jwt.JwtService
}

func New(services *services.Services, jwtService *jwt.JwtService) *Controller {
	return &Controller{services, jwtService}
}
