package services

import (
	"fmt"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

type FaqServices struct {
	repos *repositories.Repositories
}

// CreateStoreFaq Creates a new store FAQ.
func (s *FaqServices) CreateStoreFaq(input CreateStoreFaqInput) (*models.FAQ, error) {
	if !input.User.IsAdmin() && input.Store.MerchantID != input.User.ID {
		return nil, ErrUnauthorized
	}

	translation := models.Translation{
		Language: models.Language(input.Language),
		Question: input.Question,
		Answer:   input.Answer,
	}

	faq := &models.FAQ{
		Category:     *input.Category,
		StoreID:      input.Store.ID,
		IsGlobal:     false,
		Translations: []models.Translation{translation},
	}

	v := validator.New()
	if faq.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repos.Faqs.Create(faq); err != nil {
		return nil, fmt.Errorf("create store FAQ fail: %w", err)
	}

	return faq, nil
}

// CreateGlobalFaq Creates a new global FAQ
func (s *FaqServices) CreateGlobalFaq(input CreateGlobalFaqInput) (*models.FAQ, error) {
	if !input.User.IsAdmin() {
		return nil, ErrUnauthorized
	}

	translation := models.Translation{
		Language: models.Language(input.Language),
		Question: input.Question,
		Answer:   input.Answer,
	}

	faq := &models.FAQ{
		Category:        *input.Category,
		IsGlobal:        true,
		DefaultLanguage: models.Language(input.Language),
		Translations:    []models.Translation{translation},
	}

	v := validator.New()
	if faq.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repos.Faqs.Create(faq); err != nil {
		return nil, fmt.Errorf("create global FAQ fail: %w", err)
	}

	return faq, nil
}
