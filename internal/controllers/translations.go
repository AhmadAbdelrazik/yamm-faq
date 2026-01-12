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
