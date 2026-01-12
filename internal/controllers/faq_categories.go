package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/httputil"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/services"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
	"github.com/gin-gonic/gin"
)

func (c *Controller) createFaqCategoryHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)

	var input createFaqCategoryInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	category, err := c.Services.FAQCategories.Create(services.CreateCategoryInput{
		Admin:        user,
		CategoryName: input.Name,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		case errors.Is(err, services.ErrCategoryAlreadyExists):
			httputil.NewError(ctx, http.StatusConflict, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
	}

	ctx.JSON(
		http.StatusCreated,
		createFaqCategoryResponse{
			Message: fmt.Sprintf("category %v created successfully", category.Name),
		},
	)
}

func (c *Controller) getAllFaqCategoryHandler(ctx *gin.Context) {
	categories, err := c.Services.FAQCategories.GetAll()
	if err != nil {
		httputil.InternalServerError(ctx, err)
		return
	}

	fmt.Printf("categories: %v\n", categories)
	dto := make([]faqCategoryDTO, len(categories))

	for i := range dto {
		dto[i] = faqCategoryDTO(categories[i])
	}
	fmt.Printf("dto: %v\n", dto)

	ctx.JSON(http.StatusOK, getAllFaqCategoryResponse{dto})
}

func (c *Controller) updateFaqCategoryHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)

	oldCategoryName := ctx.Param("category")

	var input updateFaqCategoryInput
	if err := ctx.BindJSON(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	category, err := c.Services.FAQCategories.Update(services.UpdateCategoryInput{
		Admin:   user,
		OldName: oldCategoryName,
		NewName: input.NewName,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		case errors.Is(err, services.ErrCategoryAlreadyExists),
			errors.Is(err, services.ErrCategoryEditConflict):
			httputil.NewError(ctx, http.StatusConflict, err)
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(
		http.StatusOK,
		updateFaqCategoryResponse{
			Message: fmt.Sprintf("category %v updated successfully to %v", oldCategoryName, category.Name),
		},
	)
}

func (c *Controller) deleteFaqCategoryHandler(ctx *gin.Context) {
	user := ctx.MustGet(userContextKey).(*models.User)

	category := ctx.Param("category")

	err := c.Services.FAQCategories.Delete(services.DeleteCategoryName{
		Admin:        user,
		CategoryName: category,
	})

	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		case errors.Is(err, services.ErrCategoryNotFound):
			httputil.NotFound(ctx, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	ctx.JSON(
		http.StatusOK,
		deletedFaqCategoryResponse{
			Message: fmt.Sprintf("deleted %v category successfully", category),
		},
	)
}

type createFaqCategoryInput struct {
	Name string `json:"name"`
}

type createFaqCategoryResponse struct {
	Message string `json:"message"`
}

type getAllFaqCategoryResponse struct {
	Categories []faqCategoryDTO `json:"categories"`
}

type faqCategoryDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type updateFaqCategoryInput struct {
	NewName string `json:"new_name"`
}

type updateFaqCategoryResponse struct {
	Message string `json:"message"`
}

type deletedFaqCategoryResponse struct {
	Message string `json:"message"`
}
