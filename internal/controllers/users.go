package controllers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/httputil"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/services"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
	"github.com/gin-gonic/gin"
)

// @Summary Signup a new customer
// @Description Create a new customer account and return a session cookie
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body signupCustomerInput true "Customer Signup Details"
// @Success 201 {object} customerSignupResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 409 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /signup/customer [post]
func (c *Controller) signupCustomerHandler(ctx *gin.Context) {
	var input signupCustomerInput

	if err := ctx.BindJSON(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	user, err := c.Services.Users.SignupCustomer(services.SignupCustomerInput(input))
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrUserAlreadyExists):
			httputil.NewError(ctx, http.StatusConflict, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	c.addSessionCookie(ctx, user)

	ctx.JSON(http.StatusCreated, customerSignupResponse{
		Message:    "customer created successfully",
		CustomerID: user.ID,
	})
}

// @Summary Signup a new merchant
// @Description Create a new merchant account along with a store and return a session cookie
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body signupMerchantInput true "Merchant Signup Details"
// @Success 201 {object} merchantSignupResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 409 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /signup/merchant [post]
func (c *Controller) signupMerchantHandler(ctx *gin.Context) {
	var input signupMerchantInput

	if err := ctx.BindJSON(&input); err != nil {
		slog.Debug(err.Error())
		httputil.BadRequest(ctx, err)
		return
	}

	user, store, err := c.Services.Users.SignupMerchant(services.SignupMerchantInput(input))
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrUserAlreadyExists):
			httputil.NewError(ctx, http.StatusConflict, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	c.addSessionCookie(ctx, user)

	ctx.JSON(http.StatusCreated, merchantSignupResponse{
		Message:    "merchant created successfully",
		MerchantID: user.ID,
		StoreID:    store.ID,
	})
}

// @Summary Login user
// @Description Authenticate user and return a session cookie
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body loginInput true "Login Credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} httputil.HTTPError
// @Failure 403 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /login [post]
func (c *Controller) loginHandler(ctx *gin.Context) {
	var input loginInput

	if err := ctx.BindJSON(&input); err != nil {
		httputil.BadRequest(ctx, err)
		return
	}

	user, err := c.Services.Users.Login(services.LoginInput(input))
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrInvalid):
			httputil.InvalidEntity(ctx, err)
		case errors.Is(err, services.ErrUnauthorized):
			httputil.NewError(ctx, http.StatusForbidden, err)
		default:
			httputil.InternalServerError(ctx, err)
		}
		return
	}

	c.addSessionCookie(ctx, user)

	ctx.JSON(http.StatusOK, loginResponse{"logged in successfully"})
}

type signupCustomerInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupMerchantInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	StoreName string `json:"store_name"`
}

type customerSignupResponse struct {
	Message    string `json:"message"`
	CustomerID int    `json:"customer_id"`
}

type merchantSignupResponse struct {
	Message    string `json:"message"`
	MerchantID int    `json:"merchant_id"`
	StoreID    int    `json:"store_id"`
}

type loginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Message string `json:"message"`
}
