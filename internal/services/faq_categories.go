package services

import (
	"errors"
	"fmt"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryEditConflict  = errors.New("category edit failed")
)

type FaqCategoryService struct {
	repos *repositories.Repositories
}

// Create cretes new FAQ category. This operation can be performed by admins
// only.
func (s *FaqCategoryService) Create(input CreateCategoryInput) (*models.FAQCategory, error) {
	if !input.Admin.IsAdmin() {
		return nil, ErrUnauthorized
	}

	category := &models.FAQCategory{
		Name: input.CategoryName,
	}

	v := validator.New()
	if category.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repos.FaqCategories.Create(category); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, ErrCategoryAlreadyExists
		default:
			return nil, fmt.Errorf("create category failed: %w", err)
		}
	}

	return category, nil
}

func (s *FaqCategoryService) GetAll() ([]models.FAQCategory, error) {
	return s.repos.FaqCategories.GetAll()
}

func (s *FaqCategoryService) Find(categoryName string) (*models.FAQCategory, error) {
	category, err := s.repos.FaqCategories.FindByCategoryName(categoryName)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrCategoryNotFound
		default:
			return nil, fmt.Errorf("create category failed: %w", err)
		}
	}

	return category, nil
}

// Update updates FAQ category. This operation can be performed by admins only.
func (s *FaqCategoryService) Update(input UpdateCategoryInput) (*models.FAQCategory, error) {
	if !input.Admin.IsAdmin() {
		return nil, ErrUnauthorized
	}

	category, err := s.repos.FaqCategories.FindByCategoryName(input.OldName)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrCategoryNotFound
		default:
			return nil, fmt.Errorf("create category failed: %w", err)
		}
	}

	category.Name = input.NewName
	v := validator.New()
	if category.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repos.FaqCategories.Update(category); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, ErrCategoryAlreadyExists
		case errors.Is(err, repositories.ErrEditConflict):
			return nil, ErrCategoryEditConflict
		default:
			return nil, fmt.Errorf("category edit failed: %w", err)
		}
	}

	return category, nil
}

// Delete Hard Deletes categories. This operation can be performed by admins
// only.
func (s *FaqCategoryService) Delete(input DeleteCategoryName) error {
	if !input.Admin.IsAdmin() {
		return ErrUnauthorized
	}

	if err := s.repos.FaqCategories.HardDelete(input.CategoryName); err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return ErrCategoryNotFound
		default:
			return err
		}
	}

	return nil
}

type CreateCategoryInput struct {
	Admin        *models.User
	CategoryName string
}

type UpdateCategoryInput struct {
	Admin   *models.User
	OldName string
	NewName string
}

type DeleteCategoryName struct {
	Admin        *models.User
	CategoryName string
}
