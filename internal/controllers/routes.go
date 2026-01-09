package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes add routes to the gin engine.
func (c *Controller) Routes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// users
	api.POST("/signup/customer", c.signupCustomerHandler)
	api.POST("/signup/merchant", c.signupMerchantHandler)
	api.POST("/login", c.loginHandler)

	api.GET("/health", c.healthCheckHandler)
}

func (c *Controller) healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, healthCheckResponse{Status: "healthy"})
}

type healthCheckResponse struct {
	Status string `json:"status"`
}
