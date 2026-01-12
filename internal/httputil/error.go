package httputil

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}

	ctx.JSON(status, er)
}

func BadRequest(ctx *gin.Context, err error) {
	NewError(ctx, http.StatusBadRequest, err)
}

func NotFound(ctx *gin.Context, err error) {
	NewError(ctx, http.StatusNotFound, err)
}
func InvalidEntity(ctx *gin.Context, err error) {
	NewError(ctx, http.StatusUnprocessableEntity, err)
}

func InternalServerError(ctx *gin.Context, err error) {
	slog.Error("Internal server error", "error", err.Error())
	NewError(ctx, http.StatusInternalServerError, errors.New("something went wrong"))
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
