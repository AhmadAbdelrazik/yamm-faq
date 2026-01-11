package services

import (
	"errors"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Services struct {
	Users         *UserService
	Stores        *StoreService
	FAQCategories *FaqCategoryService
}

func New(repos *repositories.Repositories) *Services {
	return &Services{
		Users:         &UserService{repos},
		Stores:        &StoreService{repos},
		FAQCategories: &FaqCategoryService{repos},
	}
}
