package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/httputil"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/services"
	"github.com/gin-gonic/gin"
)

type contextKey string

const userContextKey contextKey = "user"

// authMiddleware add the current user to the request context to be accessed
// inside the handlers
func (c *Controller) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slog.Debug("Parsing Session Cookie")
		token, err := ctx.Cookie(jwtCookieName)
		if err != nil {
			slog.Debug("Cookie was not found with the request")
			httputil.NewError(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		claims, err := c.jwt.VerifyToken(token)
		if err != nil {
			slog.Debug(err.Error())
			httputil.NewError(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		id, err := strconv.Atoi(claims.ID)
		if err != nil {
			httputil.InternalServerError(ctx, err)
			ctx.Abort()
			return
		}

		user, err := c.Services.Users.FindByID(id)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrUnauthorized):
				slog.Debug("Attempting to access a deleted user", "userID", id)
				ctx.SetCookie("SESSION_ID", "", -1, "/", "", false, false)
				httputil.NewError(ctx, http.StatusUnauthorized, err)
			default:
				httputil.NewError(ctx, http.StatusInternalServerError, err)
			}
			ctx.Abort()
			return
		}

		slog.Debug("adding user model in the request key-value store")
		ctx.Set(userContextKey, user)

		ctx.Next()
	}
}
