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

func (c *Controller) createGlobalFaqHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)

	var input createFaqInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	category, err := c.Services.FAQCategories.Find(input.Category)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.CreateGlobalFaq(services.CreateGlobalFaqInput{
		User:     user,
		Category: category,
		Question: input.Question,
		Answer:   input.Answer,
		Language: input.Language,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
	}

	ctx.JSON(http.StatusCreated, createFaqResponse{
		Message: "FAQ Created Successfully",
		FAQ: faqDTO{
			ID:       faq.ID,
			Category: faq.Category.Name,
			Question: faq.Translations[0].Question,
			Answer:   faq.Translations[0].Answer,
			Language: string(faq.DefaultLanguage),
		},
	})
}

func (c *Controller) createStoreFaqHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)

	var input createFaqInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	store, err := c.Services.Stores.FindByMerchant(user)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrStoreNotFound):
		}
	}

	category, err := c.Services.FAQCategories.Find(input.Category)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.CreateStoreFaq(services.CreateStoreFaqInput{
		User:     user,
		Store:    store,
		Category: category,
		Question: input.Question,
		Answer:   input.Answer,
		Language: input.Language,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
	}

	ctx.JSON(http.StatusCreated, createFaqResponse{
		Message: "FAQ Created Successfully",
		FAQ: faqDTO{
			ID:       faq.ID,
			Category: faq.Category.Name,
			Question: faq.Translations[0].Question,
			Answer:   faq.Translations[0].Answer,
			Language: string(faq.DefaultLanguage),
		},
	})

}

func (c *Controller) getGlobalFaqsHandler(ctx *gin.Context) {
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

	faqs, err := c.Services.FAQs.GetGlobalFaqs(category)
	if err != nil {
		httputil.InternalServerError(ctx, err)
		return
	}

	dto := make([]faqDTO, len(faqs))

	for i := range dto {
		dto[i] = faqDTO{
			ID:       faqs[i].ID,
			Category: faqs[i].Category.Name,
			Question: faqs[i].Translations[0].Question,
			Answer:   faqs[i].Translations[0].Answer,
			Language: string(faqs[i].DefaultLanguage),
		}
	}

	ctx.JSON(http.StatusOK, getAllGlobalFaqsResponse{dto})
}

func (c *Controller) getStoreFaqsHandler(ctx *gin.Context) {
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

	faqs, err := c.Services.FAQs.GetStoreFaqs(store)
	if err != nil {
		httputil.InternalServerError(ctx, err)
		return
	}

	dto := make([]faqDTO, len(faqs))

	for i := range dto {
		dto[i] = faqDTO{
			ID:       faqs[i].ID,
			Category: faqs[i].Category.Name,
			Question: faqs[i].Translations[0].Question,
			Answer:   faqs[i].Translations[0].Answer,
			Language: string(faqs[i].DefaultLanguage),
		}
	}

	ctx.JSON(http.StatusOK, getAllGlobalFaqsResponse{dto})
}

func (c *Controller) updateGlobalFaqHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)
	faqID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	var input adminUpdateFaqInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	category, err := c.Services.FAQCategories.Find(input.Category)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	faq, err := c.Services.FAQs.AdminUpdateFaq(services.AdminUpdateFaqInput{
		Admin:    user,
		FAQID:    faqID,
		Category: category,
		IsGlobal: input.IsGlobal,
		Question: input.Question,
		Answer:   input.Answer,
		Language: input.Language,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, updateFaqResponse{
		Message: "updated FAQ successfully",
		FAQ: faqDTO{
			ID:       faq.ID,
			Category: faq.Category.Name,
			Question: faq.Translations[0].Question,
			Answer:   faq.Translations[0].Answer,
			Language: string(faq.Translations[0].Language),
		},
	})
}

func (c *Controller) updateStoreFaqHandler(ctx *gin.Context) {
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

	var input merchantUpdateFaqInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	category, err := c.Services.FAQCategories.Find(input.Category)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
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

	faq, err := c.Services.FAQs.MerchantUpdateFaq(services.MerchantUpdateFaqInput{
		Merchant: user,
		Store:    store,
		FAQID:    faqID,
		Category: category,
		Question: input.Question,
		Answer:   input.Answer,
		Language: input.Language,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, updateFaqResponse{
		Message: "updated FAQ successfully",
		FAQ: faqDTO{
			ID:       faq.ID,
			Category: faq.Category.Name,
			Question: faq.Translations[0].Question,
			Answer:   faq.Translations[0].Answer,
			Language: string(faq.Translations[0].Language),
		},
	})
}

func (c *Controller) deleteGlobalFaqHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)
	faqID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	if err := c.Services.FAQs.AdminDelete(user, faqID); err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, deleteFaqResponse{"deleted FAQ successfully"})
}

func (c *Controller) deleteStoreFaqHandler(ctx *gin.Context) {
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

	if err := c.Services.FAQs.MerchantDelete(user, store, faqID); err != nil {
		switch {
		case errors.Is(err, services.ErrFaqNotFound):
			httputil.NotFound(ctx, err)
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
	}

	ctx.JSON(http.StatusOK, deleteFaqResponse{"deleted FAQ successfully"})
}

type createFaqInput struct {
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Language string `json:"language"`
}

type adminUpdateFaqInput struct {
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Language string `json:"language"`
	IsGlobal bool   `json:"is_global"`
}

type merchantUpdateFaqInput struct {
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Language string `json:"language"`
}

type createFaqResponse struct {
	Message string `json:"message"`
	FAQ     faqDTO `json:"faq"`
}

type updateFaqResponse struct {
	Message string `json:"message"`
	FAQ     faqDTO `json:"faq"`
}

type deleteFaqResponse struct {
	Message string `json:"message"`
}

type faqDTO struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Language string `json:"language"`
}

type getAllGlobalFaqsResponse struct {
	FAQs []faqDTO `json:"faqs"`
}
