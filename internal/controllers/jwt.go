package controllers

import (
	"net/http"
	"time"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/jwt"
	"github.com/gin-gonic/gin"
)

const jwtCookieName = "JWTSESSION"

func (c *Controller) addSessionCookie(ctx *gin.Context, user *models.User) error {
	claims := jwt.NewUserClaims(user.ID, string(user.Role))
	token, err := c.jwt.GenerateToken(claims)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}

	ctx.SetCookieData(cookie)

	return nil
}
