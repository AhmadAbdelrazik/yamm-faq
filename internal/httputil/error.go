package httputil

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	er := httpError{
		Code:    status,
		Message: err.Error(),
	}

	if status == http.StatusInternalServerError {
		slog.Error("Internal server error", "error", err.Error())
		er.Message = "Something went wrong"
	}

	ctx.JSON(status, er)
}

func NewValidationError(ctx *gin.Context, errors map[string]string) {
	er := validationError{
		Code:    http.StatusBadRequest,
		Message: "invalid request body",
		Errors:  errors,
	}
	ctx.JSON(http.StatusBadRequest, er)
}

// HTTPError example
type httpError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

type validationError struct {
	Code    int               `json:"code" example:"400"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}
