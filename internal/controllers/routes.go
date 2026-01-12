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
	api.GET("/faq-categories", c.getAllFaqCategoryHandler)

	auth.POST("/faq-categories", c.createFaqCategoryHandler)
	auth.PUT("/faq-categories/:category", c.updateFaqCategoryHandler)
	auth.DELETE("/faq-categories/:category", c.deleteFaqCategoryHandler)

	// Global FAQs
	api.GET("/faq-categories/:category", c.getGlobalFaqsHandler)

	auth.POST("/faq-categories/:category", c.createGlobalFaqHandler)
	auth.PUT("/faq-categories/:category/:id", c.updateGlobalFaqHandler)
	auth.DELETE("/faq-categories/:category/:id", c.deleteGlobalFaqHandler)

	// Store specific FAQs
	api.GET("/stores/:id/faqs", c.getStoreFaqsHandler)

	auth.POST("/stores/:id/faqs", c.createStoreFaqHandler)
	auth.PUT("/stores/:id/faqs/:faq-id", c.updateStoreFaqHandler)
	auth.DELETE("/stores/:id/faqs/:faq-id", c.deleteStoreFaqHandler)

	// Global FAQs Translations
	api.GET("/faq-categories/:category/:id/translations", c.getGlobalFaqTranslationsHandler)
	api.GET("/faq-categories/:category/:id/:language", c.getGlobalFaqLanguageHandler)

	auth.POST("/faq-categories/:category/:id/translations", c.createGlobalFaqTranslationHandler)
	auth.PUT("/faq-categories/:category/:id/:language", c.updateGlobalFaqLanguageHandler)
	auth.DELETE("/faq-categories/:category/:id/:language", c.deleteGlobalFaqLanguageHandler)

	// Store Specific FAQs Translations
	api.GET("/stores/:id/faqs/:faq-id/translations", c.getStoreFaqTranslationsHandler)
	api.GET("/stores/:id/faqs/:faq-id/:language", c.getStoreFaqLanguageHandler)

	auth.POST("/stores/:id/faqs/:faq-id/translations", c.createStoreFaqTranslationHandler)
	auth.PUT("/stores/:id/faqs/:faq-id/:language", c.updateStoreFaqLanguageHandler)
	auth.DELETE("/stores/:id/faqs/:faq-id/:language", c.deleteStoreFaqLanguageHandler)

	api.GET("/health", c.healthCheckHandler)
}

func (c *Controller) healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, healthCheckResponse{Status: "healthy"})
}

type healthCheckResponse struct {
	Status string `json:"status"`
}
