package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes add routes to the gin engine.
func (c *Controller) Routes(r *gin.Engine) {
	api := r.Group("/api/v1")

	auth := api.Group("/")
	auth.Use(c.authMiddleware())

	// users
	api.POST("/signup/customer", c.signupCustomerHandler)
	api.POST("/signup/merchant", c.signupMerchantHandler)
	api.POST("/login", c.loginHandler)

	// FAQ Categories
	api.GET("/faq-categories/:category", c.getFaqCategoryHandler)
	api.GET("/faq-categories", c.getAllFaqCategoryHandler)

	auth.POST("/faq-categories", c.createFaqCategoryHandler)
	auth.PUT("/faq-categories/:category", c.updateFaqCategoryHandler)
	auth.DELETE("/faq-categories/:category", c.deleteFaqCategoryHandler)

	api.GET("/health", c.healthCheckHandler)
}

func (c *Controller) healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, healthCheckResponse{Status: "healthy"})
}

type healthCheckResponse struct {
	Status string `json:"status"`
}
