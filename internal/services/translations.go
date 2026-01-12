package services

import (
	"errors"
	"fmt"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

var (
	ErrTranslationNotFound        = errors.New("translation not found")
	ErrTranslationAlreadyExists   = errors.New("translation already exists")
	ErrDeletingDefaultTranslation = errors.New("can't delete default translation")
)

type TranslationService struct {
	repo *repositories.Repositories
}

func (s *TranslationService) GetTranslation(faqID int, language models.Language) (*models.Translation, error) {
	v := validator.New()
	if language.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	translation, err := s.repo.Translations.Get(faqID, language)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrTranslationNotFound
		default:
			return nil, err
		}
	}

	return translation, nil
}

// GlobalCreate Create Translation on Any FAQ. Allowed for Admins only
func (s *TranslationService) GlobalCreate(input CreateGlobalTranslationInput) (*models.Translation, error) {
	if !input.User.IsAdmin() {
		return nil, ErrUnauthorized
	}

	translation := &models.Translation{
		FAQID:    input.FAQ.ID,
		Language: models.Language(input.Language),
		Question: input.Question,
		Answer:   input.Answer,
	}

	v := validator.New()
	if translation.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repo.Translations.Create(translation); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, ErrTranslationAlreadyExists
		default:
			return nil, err
		}
	}

	return translation, nil
}

// StoreCreate Create FAQ Translation related to a store. Can be performed by
// admins or store owners
func (s *TranslationService) StoreCreate(input CreateStoreTranslationInput) (*models.Translation, error) {
	if !input.User.IsAdmin() && input.User.ID != input.Store.MerchantID {
		return nil, ErrUnauthorized
	}

	if input.FAQ.StoreID != input.Store.ID {
		return nil, fmt.Errorf("%w: faq doesn't belong to this store", ErrUnauthorized)
	}

	translation := &models.Translation{
		FAQID:    input.FAQ.ID,
		Language: models.Language(input.Language),
		Question: input.Question,
		Answer:   input.Answer,
	}

	v := validator.New()
	if translation.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repo.Translations.Create(translation); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, ErrTranslationAlreadyExists
		default:
			return nil, err
		}
	}

	return translation, nil
}

// GlobalUpdate Update Translations globally. Allowed by Admins Only
func (s *TranslationService) GlobalUpdate(input UpdateGlobalTranslationInput) (*models.Translation, error) {
	if !input.User.IsAdmin() {
		return nil, ErrUnauthorized
	}

	translation, err := s.repo.Translations.Get(input.FAQ.ID, models.Language(input.CurrentLanguage))
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrTranslationNotFound
		default:
			return nil, err
		}
	}

	translation.Language = models.Language(input.NewLanguage)
	translation.Question = input.Question
	translation.Answer = input.Answer

	v := validator.New()
	if translation.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repo.Translations.Update(translation); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, ErrTranslationAlreadyExists
		default:
			return nil, err
		}
	}

	return translation, nil
}

// StoreUpdate Update FAQ Translation related to a store. Can be performed by
// admins or store owners
func (s *TranslationService) StoreUpdate(input UpdateStoreTranslationInput) (*models.Translation, error) {
	if !input.User.IsAdmin() && input.User.ID != input.Store.MerchantID {
		return nil, ErrUnauthorized
	}

	if input.FAQ.StoreID != input.Store.ID {
		return nil, fmt.Errorf("%w: faq doesn't belong to this store", ErrUnauthorized)
	}

	translation, err := s.repo.Translations.Get(input.FAQ.ID, models.Language(input.CurrentLanguage))
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrTranslationNotFound
		default:
			return nil, err
		}
	}

	translation.Language = models.Language(input.NewLanguage)
	translation.Question = input.Question
	translation.Answer = input.Answer

	v := validator.New()
	if translation.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repo.Translations.Create(translation); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, ErrTranslationAlreadyExists
		default:
			return nil, err
		}
	}

	return translation, nil
}

// GlobalDelete Deletes translations globally. Allowed by admins only
func (s *TranslationService) GlobalDelete(input DeleteGlobalTranslationInput) error {
	if !input.User.IsAdmin() {
		return ErrUnauthorized
	}

	if input.FAQ.DefaultLanguage == models.Language(input.Language) {
		return ErrDeletingDefaultTranslation
	}

	fmt.Printf("input.FAQ: %#v\n", input.FAQ)
	fmt.Printf("input.Language: %v\n", input.Language)

	err := s.repo.Translations.DeleteByLanguage(input.FAQ.ID, input.Language)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return ErrTranslationNotFound
		default:
			return err
		}
	}

	return nil
}

// StoreDelete Delete FAQ Translation related to a store. Can be performed by
// admins or store owners
func (s *TranslationService) StoreDelete(input DeleteStoreTranslationInput) error {
	if !input.User.IsAdmin() && input.User.ID != input.Store.MerchantID {
		return ErrUnauthorized
	}

	if input.FAQ.StoreID != input.Store.ID {
		return fmt.Errorf("%w: faq doesn't belong to this store", ErrUnauthorized)
	}

	if input.FAQ.DefaultLanguage == models.Language(input.Language) {
		return ErrDeletingDefaultTranslation
	}

	err := s.repo.Translations.DeleteByLanguage(input.FAQ.ID, input.Language)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return ErrTranslationNotFound
		default:
			return err
		}
	}

	return nil
}

type CreateGlobalTranslationInput struct {
	User     *models.User
	FAQ      *models.FAQ
	Language string
	Question string
	Answer   string
}

type CreateStoreTranslationInput struct {
	User     *models.User
	Store    *models.Store
	FAQ      *models.FAQ
	Language string
	Question string
	Answer   string
}

type UpdateGlobalTranslationInput struct {
	User            *models.User
	FAQ             *models.FAQ
	CurrentLanguage string
	NewLanguage     string
	Question        string
	Answer          string
}

type UpdateStoreTranslationInput struct {
	User            *models.User
	Store           *models.Store
	FAQ             *models.FAQ
	CurrentLanguage string
	NewLanguage     string
	Question        string
	Answer          string
}

type DeleteGlobalTranslationInput struct {
	User     *models.User
	FAQ      *models.FAQ
	Language string
}

type DeleteStoreTranslationInput struct {
	User     *models.User
	Store    *models.Store
	FAQ      *models.FAQ
	Language string
}
