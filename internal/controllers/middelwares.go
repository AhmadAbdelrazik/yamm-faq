package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/httputil"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
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

// rateLimitMiddleware limits the number of requests for each user
func rateLimitMiddleware(limitRate float64, burst int, cleanupDuration time.Duration) gin.HandlerFunc {
	type Limiter struct {
		limit       *rate.Limiter
		lastRequest time.Time
	}

	freq := struct {
		m map[string]*Limiter
		sync.RWMutex
	}{
		m: make(map[string]*Limiter),
	}

	// cleanup
	go func() {
		freq.RWMutex.Lock()
		for ip, limiter := range freq.m {
			if time.Since(limiter.lastRequest) > cleanupDuration {
				delete(freq.m, ip)
			}
		}
		freq.RWMutex.Unlock()
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		freq.RWMutex.RLock()
		_, ok := freq.m[ip]
		freq.RWMutex.RUnlock()

		freq.RWMutex.Lock()
		if !ok {
			freq.m[ip] = &Limiter{
				limit:       rate.NewLimiter(rate.Limit(limitRate), burst),
				lastRequest: time.Now(),
			}
		} else {
			if !freq.m[ip].limit.Allow() {
				httputil.NewError(c, http.StatusTooManyRequests, errors.New("Too many requests"))
				c.Abort()
			}
			freq.m[ip].lastRequest = time.Now()
		}
		freq.RWMutex.Unlock()

		c.Next()
	}
}
