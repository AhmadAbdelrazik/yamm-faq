package services

import (
	"errors"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Services struct {
	Users *UserService
}

func New(repos *repositories.Repositories) *Services {
	return &Services{
		Users: &UserService{repos},
	}
}
