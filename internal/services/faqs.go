package services

import (
	"errors"
	"fmt"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

var (
	ErrFaqNotFound = errors.New("faq not found")
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
		Category:        *input.Category,
		StoreID:         input.Store.ID,
		IsGlobal:        false,
		Translations:    []models.Translation{translation},
		DefaultLanguage: models.Language(input.Language),
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

// FindFAQInCategory Finds FAQ with all its available translations in a category
func (s *FaqServices) FindFAQInCategory(faqID int, category *models.FAQCategory) (*models.FAQ, error) {
	faq, err := s.repos.Faqs.Find(faqID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrFaqNotFound
		default:
			return nil, fmt.Errorf("Find FAQ Fail: %w", err)
		}
	}

	if faq.Category.Name != category.Name {
		return nil, ErrFaqNotFound
	}

	return faq, nil
}

// FindFAQInStore Finds FAQ with all its available translations in a store
func (s *FaqServices) FindFAQInStore(faqID int, store *models.Store) (*models.FAQ, error) {
	faq, err := s.repos.Faqs.Find(faqID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrFaqNotFound
		default:
			return nil, fmt.Errorf("Find FAQ Fail: %w", err)
		}
	}

	if faq.StoreID != store.ID {
		return nil, ErrFaqNotFound
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

func (s *FaqServices) GetGlobalFaqs(category *models.FAQCategory) ([]models.FAQ, error) {
	return s.repos.Faqs.GetAllByCategory(category.Name)
}

func (s *FaqServices) GetStoreFaqs(store *models.Store) ([]models.FAQ, error) {
	return s.repos.Faqs.GetAllByStore(store.ID)
}

// MerchantUpdateFaq updates FAQ related to merchant's store.
func (s *FaqServices) MerchantUpdateFaq(input MerchantUpdateFaqInput) (*models.FAQ, error) {
	if input.Merchant.ID != input.Store.MerchantID {
		return nil, fmt.Errorf("%w: user is not the owner of the store", ErrUnauthorized)
	}

	faq, err := s.repos.Faqs.FindDefault(input.FAQID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrFaqNotFound
		default:
			return nil, err
		}
	}

	fmt.Printf("faq: %v\n", faq)
	fmt.Printf("input.Store: %v\n", input.Store)

	if input.Store.ID != faq.StoreID {
		return nil, fmt.Errorf("%w: faq doesn't belong to the store", ErrUnauthorized)
	}

	faq.Translations[0].Question = input.Question
	faq.Translations[0].Answer = input.Answer
	faq.Translations[0].Language = models.Language(input.Language)

	faq.Category = *input.Category
	faq.DefaultLanguage = models.Language(input.Language)

	v := validator.New()
	if faq.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repos.Faqs.Update(faq); err != nil {
		return nil, err
	}
	if err := s.repos.Translations.Update(&faq.Translations[0]); err != nil {
		return nil, err
	}

	return faq, nil
}

// AdminUpdateFaq updates FAQ related to merchant's store.
func (s *FaqServices) AdminUpdateFaq(input AdminUpdateFaqInput) (*models.FAQ, error) {
	if !input.Admin.IsAdmin() {
		return nil, ErrUnauthorized
	}

	faq, err := s.repos.Faqs.FindDefault(input.FAQID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrFaqNotFound
		default:
			return nil, err
		}
	}

	faq.Translations[0].Question = input.Question
	faq.Translations[0].Answer = input.Answer
	faq.Translations[0].Language = models.Language(input.Language)

	faq.Category = *input.Category
	faq.DefaultLanguage = models.Language(input.Language)
	faq.IsGlobal = input.IsGlobal

	v := validator.New()
	if faq.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repos.Faqs.Update(faq); err != nil {
		return nil, err
	}
	if err := s.repos.Translations.Update(&faq.Translations[0]); err != nil {
		return nil, err
	}

	return faq, nil
}

func (s *FaqServices) AdminDelete(user *models.User, faqID int) error {
	if !user.IsAdmin() {
		return ErrUnauthorized
	}

	err := s.repos.Faqs.Delete(faqID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return ErrFaqNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *FaqServices) MerchantDelete(merchant *models.User, store *models.Store, faqID int) error {
	if merchant.ID != store.MerchantID {
		return ErrUnauthorized
	}

	err := s.repos.Faqs.MerchantDelete(faqID, store.ID)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return ErrUnauthorized
		default:
			return err
		}
	}

	return nil
}

type CreateStoreFaqInput struct {
	User     *models.User
	Store    *models.Store
	Category *models.FAQCategory
	Question string
	Answer   string
	Language string
}

type CreateGlobalFaqInput struct {
	User     *models.User
	Category *models.FAQCategory
	Question string
	Answer   string
	Language string
}

type MerchantUpdateFaqInput struct {
	Merchant *models.User
	Store    *models.Store
	FAQID    int
	Category *models.FAQCategory
	Question string
	Answer   string
	Language string
}

type AdminUpdateFaqInput struct {
	Admin    *models.User
	FAQID    int
	Category *models.FAQCategory
	IsGlobal bool
	Question string
	Answer   string
	Language string
}
