package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/httputil"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/services"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
	"github.com/gin-gonic/gin"
)

// @Summary      Get all translations for a global FAQ
// @Description  Retrieve every available language translation for a specific global FAQ
// @Tags         Global FAQ Translations
// @Produce      json
// @Param        category  path      string  true  "Category Name"
// @Param        id        path      int     true  "FAQ ID"
// @Success      200       {object}  getFaqTranslationsResponse
// @Failure      400       {object}  httputil.HTTPError
// @Failure      404       {object}  httputil.HTTPError "Category or FAQ not found"
// @Failure      500       {object}  httputil.HTTPError
// @Router       /faq-categories/{category}/{id}/translations [get]
func (c *Controller) getGlobalFaqTranslationsHandler(ctx *gin.Context) {
	categoryName := ctx.Param("category")
	faqID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	category, err := c.Services.FAQCategories.Find(categoryName)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.FindFAQInCategory(faqID, category)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	translations := make([]translationDTO, len(faq.Translations))

	for i, t := range faq.Translations {
		translations[i] = translationDTO{
			ID:       t.ID,
			FAQID:    t.FAQID,
			Language: string(t.Language),
			Question: t.Question,
			Answer:   t.Answer,
		}
	}

	ctx.JSON(http.StatusOK, getFaqTranslationsResponse{
		FaqID:           faqID,
		Category:        faq.Category.Name,
		DefaultLanguage: string(faq.DefaultLanguage),
		Translations:    translations,
	})

}

// @Summary      Get global FAQ translation by language
// @Description  Retrieve a specific language version of a global FAQ
// @Tags         Global FAQ Translations
// @Produce      json
// @Param        category  path      string  true  "Category Name"
// @Param        id        path      int     true  "FAQ ID"
// @Param        language  path      string  true  "Language Code (e.g., 'en', 'ar')"
// @Success      200       {object}  translationDTO
// @Failure      404       {object}  httputil.HTTPError "Translation not found"
// @Failure      422       {object}  httputil.HTTPError "Invalid language code"
// @Router       /faq-categories/{category}/{id}/{language} [get]
func (c *Controller) getGlobalFaqLanguageHandler(ctx *gin.Context) {
	faqID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}
	language := ctx.Param("language")

	translation, err := c.Services.Translations.GetTranslation(faqID, models.Language(language))
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrTranslationNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, translationDTO{
		ID:       translation.ID,
		FAQID:    faqID,
		Language: string(translation.Language),
		Question: translation.Question,
		Answer:   translation.Answer,
	})

}

// @Summary      Create global FAQ translation
// @Description  Add a new language translation to an existing global FAQ (Admin only)
// @Tags         Global FAQ Translations
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        category  path      string                  true  "Category Name"
// @Param        id        path      int                     true  "FAQ ID"
// @Param        input     body      createTranslationInput  true  "Translation Content"
// @Success      201       {object}  createTranslationResponse
// @Failure      403       {object}  httputil.HTTPError "Unauthorized"
// @Failure      409       {object}  httputil.HTTPError "Translation already exists"
// @Failure      422       {object}  httputil.HTTPError "Validation failed"
// @Router       /faq-categories/{category}/{id}/translations [post]
func (c *Controller) createGlobalFaqTranslationHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)
	faqID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	var input createTranslationInput
	if err := ctx.Bind(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	categoryName := ctx.Param("category")
	category, err := c.Services.FAQCategories.Find(categoryName)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.FindFAQInCategory(faqID, category)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	translation, err := c.Services.Translations.GlobalCreate(services.CreateGlobalTranslationInput{
		User:     user,
		FAQ:      faq,
		Language: input.Language,
		Question: input.Question,
		Answer:   input.Answer,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusConflict, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrTranslationAlreadyExists):
			httputil.NewError(ctx, http.StatusConflict, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusCreated, createTranslationResponse{
		Message: "translation added successfully",
		Translation: translationDTO{
			ID:       translation.ID,
			FAQID:    translation.FAQID,
			Language: string(translation.Language),
			Question: translation.Question,
			Answer:   translation.Answer,
		},
	})
}

// @Summary      Update global FAQ translation
// @Description  Modify an existing translation for a global FAQ (Admin only)
// @Tags         Global FAQ Translations
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        category  path      string                  true  "Category Name"
// @Param        id        path      int                     true  "FAQ ID"
// @Param        language  path      string                  true  "Current Language"
// @Param        input     body      updateTranslationInput  true  "Updated Content"
// @Success      200       {object}  updateTranslationResponse
// @Failure      404       {object}  httputil.HTTPError "Translation not found"
// @Failure      409       {object}  httputil.HTTPError "Unauthorized or Conflict"
// @Router       /faq-categories/{category}/{id}/{language} [put]
func (c *Controller) updateGlobalFaqLanguageHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)
	language := ctx.Param("language")
	faqID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	var input updateTranslationInput
	if err := ctx.Bind(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	categoryName := ctx.Param("category")
	category, err := c.Services.FAQCategories.Find(categoryName)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.FindFAQInCategory(faqID, category)
	if err != nil {
		switch {

		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	translation, err := c.Services.Translations.GlobalUpdate(services.UpdateGlobalTranslationInput{
		User:            user,
		FAQ:             faq,
		CurrentLanguage: language,
		NewLanguage:     input.Language,
		Question:        input.Question,
		Answer:          input.Answer,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusConflict, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrTranslationNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, updateTranslationResponse{
		Message: "updated Translation successfully",
		Translation: translationDTO{
			ID:       translation.ID,
			FAQID:    translation.FAQID,
			Language: string(translation.Language),
			Question: translation.Question,
			Answer:   translation.Answer,
		},
	})
}

// @Summary      Delete global FAQ translation
// @Description  Remove a specific language translation (Cannot delete default language) (Admin only)
// @Tags         Global FAQ Translations
// @Security     ApiKeyAuth
// @Param        category  path      string  true  "Category Name"
// @Param        id        path      int     true  "FAQ ID"
// @Param        language  path      string  true  "Language to Delete"
// @Success      200       {object}  deleteTranslationResponse
// @Failure      409       {object}  httputil.HTTPError "Cannot delete default language or Unauthorized"
// @Router       /faq-categories/{category}/{id}/{language} [delete]
func (c *Controller) deleteGlobalFaqLanguageHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)
	language := ctx.Param("language")
	faqID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	categoryName := ctx.Param("category")
	category, err := c.Services.FAQCategories.Find(categoryName)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}
	faq, err := c.Services.FAQs.FindFAQInCategory(faqID, category)
	if err != nil {
		switch {

		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	err = c.Services.Translations.GlobalDelete(services.DeleteGlobalTranslationInput{
		User:     user,
		FAQ:      faq,
		Language: language,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized),
			errors.Is(err, services.ErrDeletingDefaultTranslation):
			httputil.NewError(ctx, http.StatusConflict, err)
		case errors.Is(err, services.ErrTranslationNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, deleteTranslationResponse{
		Message: "Translation has been deleted successfully",
	})
}

// @Summary      Get all translations for a store FAQ
// @Description  Retrieve every available language translation for a specific store's FAQ
// @Tags         Store FAQ Translations
// @Produce      json
// @Param        id      path      int  true  "Store ID"
// @Param        faq-id  path      int  true  "FAQ ID"
// @Success      200     {object}  getFaqTranslationsResponse
// @Failure      404     {object}  httputil.HTTPError "Store or FAQ not found"
// @Router       /stores/{id}/faqs/{faq-id}/translations [get]
func (c *Controller) getStoreFaqTranslationsHandler(ctx *gin.Context) {
	faqID, err := strconv.Atoi(ctx.Param("faq-id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	storeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}
	store, err := c.Services.Stores.FindByID(storeID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrStoreNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.FindFAQInStore(faqID, store)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	translations := make([]translationDTO, len(faq.Translations))

	for i, t := range faq.Translations {
		translations[i] = translationDTO{
			ID:       t.ID,
			FAQID:    t.FAQID,
			Language: string(t.Language),
			Question: t.Question,
			Answer:   t.Answer,
		}
	}

	ctx.JSON(http.StatusOK, getFaqTranslationsResponse{
		FaqID:           faqID,
		Category:        faq.Category.Name,
		DefaultLanguage: string(faq.DefaultLanguage),
		Translations:    translations,
	})
}

// @Summary      Get store FAQ translation by language
// @Description  Retrieve a specific language version of a store-specific FAQ
// @Tags         Store FAQ Translations
// @Produce      json
// @Param        id        path      int     true  "Store ID"
// @Param        faq-id    path      int     true  "FAQ ID"
// @Param        language  path      string  true  "Language Code"
// @Success      200       {object}  translationDTO
// @Router       /stores/{id}/faqs/{faq-id}/{language} [get]
func (c *Controller) getStoreFaqLanguageHandler(ctx *gin.Context) {
	faqID, err := strconv.Atoi(ctx.Param("faq-id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}
	language := ctx.Param("language")

	translation, err := c.Services.Translations.GetTranslation(faqID, models.Language(language))
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrTranslationNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, translationDTO{
		ID:       translation.ID,
		FAQID:    faqID,
		Language: string(translation.Language),
		Question: translation.Question,
		Answer:   translation.Answer,
	})
}

// @Summary      Create store FAQ translation
// @Description  Add a new language translation to an existing store FAQ (Merchant only)
// @Tags         Store FAQ Translations
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id      path      int                     true  "Store ID"
// @Param        faq-id  path      int                     true  "FAQ ID"
// @Param        input   body      createTranslationInput  true  "Translation Content"
// @Success      201     {object}  createTranslationResponse
// @Failure      409     {object}  httputil.HTTPError "Already exists or Unauthorized"
// @Router       /stores/{id}/faqs/{faq-id}/translations [post]
func (c *Controller) createStoreFaqTranslationHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)
	storeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}
	faqID, err := strconv.Atoi(ctx.Param("faq-id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	var input createTranslationInput
	if err := ctx.Bind(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	store, err := c.Services.Stores.FindByID(storeID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrStoreNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.FindFAQInStore(faqID, store)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	translation, err := c.Services.Translations.StoreCreate(services.CreateStoreTranslationInput{
		User:     user,
		Store:    store,
		FAQ:      faq,
		Language: input.Language,
		Question: input.Question,
		Answer:   input.Answer,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusConflict, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrTranslationAlreadyExists):
			httputil.NewError(ctx, http.StatusConflict, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusCreated, createTranslationResponse{
		Message: "translation added successfully",
		Translation: translationDTO{
			ID:       translation.ID,
			FAQID:    translation.FAQID,
			Language: string(translation.Language),
			Question: translation.Question,
			Answer:   translation.Answer,
		},
	})
}

// @Summary      Update store FAQ translation
// @Description  Modify an existing translation for a store FAQ (Merchant only)
// @Tags         Store FAQ Translations
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id        path      int                     true  "Store ID"
// @Param        faq-id    path      int                     true  "FAQ ID"
// @Param        language  path      string                  true  "Current Language"
// @Param        input     body      updateTranslationInput  true  "Updated Content"
// @Success      200       {object}  updateTranslationResponse
// @Router       /stores/{id}/faqs/{faq-id}/{language} [put]
func (c *Controller) updateStoreFaqLanguageHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)
	language := ctx.Param("language")
	storeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}
	faqID, err := strconv.Atoi(ctx.Param("faq-id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	var input updateTranslationInput
	if err := ctx.Bind(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	store, err := c.Services.Stores.FindByID(storeID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrStoreNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.FindFAQInStore(faqID, store)
	if err != nil {
		switch {

		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	translation, err := c.Services.Translations.StoreUpdate(services.UpdateStoreTranslationInput{
		User:            user,
		Store:           store,
		FAQ:             faq,
		CurrentLanguage: language,
		NewLanguage:     input.Language,
		Question:        input.Question,
		Answer:          input.Answer,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusConflict, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrTranslationNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, updateTranslationResponse{
		Message: "updated Translation successfully",
		Translation: translationDTO{
			ID:       translation.ID,
			FAQID:    translation.FAQID,
			Language: string(translation.Language),
			Question: translation.Question,
			Answer:   translation.Answer,
		},
	})
}

// @Summary      Delete store FAQ translation
// @Description  Remove a translation from a store FAQ (Merchant only)
// @Tags         Store FAQ Translations
// @Security     ApiKeyAuth
// @Param        id        path      int     true  "Store ID"
// @Param        faq-id    path      int     true  "FAQ ID"
// @Param        language  path      string  true  "Language to Delete"
// @Success      200       {object}  deleteTranslationResponse
// @Failure      409       {object}  httputil.HTTPError "Cannot delete default language"
// @Router       /stores/{id}/faqs/{faq-id}/{language} [delete]
func (c *Controller) deleteStoreFaqLanguageHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)
	language := ctx.Param("language")
	storeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}
	faqID, err := strconv.Atoi(ctx.Param("faq-id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	var input createTranslationInput
	if err := ctx.Bind(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	store, err := c.Services.Stores.FindByID(storeID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrStoreNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	categoryName := ctx.Param("category")
	category, err := c.Services.FAQCategories.Find(categoryName)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}
	faq, err := c.Services.FAQs.FindFAQInCategory(faqID, category)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	err = c.Services.Translations.StoreDelete(services.DeleteStoreTranslationInput{
		User:     user,
		FAQ:      faq,
		Language: language,
		Store:    store,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized),
			errors.Is(err, services.ErrDeletingDefaultTranslation):
			httputil.NewError(ctx, http.StatusConflict, err)
		case errors.Is(err, services.ErrTranslationNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, deleteTranslationResponse{
		Message: "Translation has been deleted successfully",
	})

}

type translationDTO struct {
	ID       int    `json:"id"`
	FAQID    int    `json:"faq_id"`
	Language string `json:"language"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type getFaqTranslationsResponse struct {
	FaqID           int              `json:"id"`
	Category        string           `json:"category"`
	DefaultLanguage string           `json:"default_language"`
	Translations    []translationDTO `json:"translations"`
}

type createTranslationInput struct {
	Language string `json:"language"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type updateTranslationInput struct {
	Language string `json:"language"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type createTranslationResponse struct {
	Message     string         `json:"message"`
	Translation translationDTO `json:"translation"`
}

type updateTranslationResponse struct {
	Message     string         `json:"message"`
	Translation translationDTO `json:"translation"`
}

type deleteTranslationResponse struct {
	Message string `json:"message"`
}
