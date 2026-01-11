package services

import (
	"errors"
	"fmt"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

var (
	ErrStoreNotFound = errors.New("store not found")
)

type StoreService struct {
	repos *repositories.Repositories
}

func (s *StoreService) FindByID(storeID int) (*models.Store, error) {
	store, err := s.repos.Stores.FindByID(storeID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrStoreNotFound
		default:
			return nil, fmt.Errorf("store find by id fail: %w", err)
		}
	}

	return store, nil
}

func (s *StoreService) FindByMerchant(merchant *models.User) (*models.Store, error) {
	if !merchant.IsMerchant() {
		return nil, ErrUnauthorized
	}

	store, err := s.repos.Stores.FindByMerchantID(merchant.ID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrStoreNotFound
		default:
			return nil, fmt.Errorf("store find by id fail: %w", err)
		}
	}

	return store, nil
}

func (s *StoreService) Update(input UpdateStoreInput) error {
	store, err := s.repos.Stores.FindByMerchantID(input.merchant.ID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return ErrStoreNotFound
		default:
			return fmt.Errorf("store update fail: %w", err)
		}
	}

	if input.merchant.ID != store.MerchantID {
		return ErrUnauthorized
	}

	store.Name = input.NewName

	v := validator.New()
	if store.Validate(v); !v.Valid() {
		return v.Err()
	}

	if err := s.repos.Stores.Update(store); err != nil {
		return fmt.Errorf("store update fail: %w", err)
	}

	return nil
}

type UpdateStoreInput struct {
	merchant *models.Store
	NewName  string
}
