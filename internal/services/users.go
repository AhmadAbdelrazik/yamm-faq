package services

import (
	"errors"
	"fmt"

	"github.com/AhmadAbdelrazik/yamm_faq/internal/models"
	"github.com/AhmadAbdelrazik/yamm_faq/internal/repositories"
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserService struct {
	repos *repositories.Repositories
}

// SignupCustomer creates a new customer
func (s *UserService) SignupCustomer(input SignupCustomerInput) (*models.User, error) {
	pass, _ := models.NewPassword(input.Password)

	user := &models.User{
		Email:    input.Email,
		Password: pass,
		Role:     models.RoleCustomer,
	}

	// validate that user is not an admin
	v := validator.New()
	if user.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repos.Users.Create(user); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, ErrUserAlreadyExists
		default:
			return nil, fmt.Errorf("user signup failed: %w", err)
		}
	}

	return user, nil
}

// SignupMerchant creates a new merchant
func (s *UserService) SignupMerchant(input SignupMerchantInput) (*models.User, *models.Store, error) {
	pass, _ := models.NewPassword(input.Password)

	user := &models.User{
		Email:    input.Email,
		Password: pass,
		Role:     models.RoleMerchant,
	}

	v := validator.New()
	if user.Validate(v); !v.Valid() {
		return nil, nil, v.Err()
	}

	store := &models.Store{
		Name: input.StoreName,
	}

	if store.Validate(v); !v.Valid() {
		return nil, nil, v.Err()
	}

	if err := s.repos.Users.CreateMerchant(user, store); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, nil, ErrUserAlreadyExists
		default:
			return nil, nil, fmt.Errorf("user signup failed: %w", err)
		}
	}

	return user, store, nil
}

// SignupAdmin Allows existing admin to create a new user admin.
func (s *UserService) SignupAdmin(input SignupAdminInput) (*models.User, error) {
	if !input.Admin.IsAdmin() {
		return nil, ErrUnauthorized
	}

	pass, _ := models.NewPassword(input.Password)

	user := &models.User{
		Email:    input.Email,
		Password: pass,
		Role:     models.RoleAdmin,
	}

	v := validator.New()
	if user.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	if err := s.repos.Users.Create(user); err != nil {
		switch {
		case errors.Is(err, repositories.ErrDuplicate):
			return nil, ErrUserAlreadyExists
		default:
			return nil, fmt.Errorf("user signup failed: %w", err)
		}
	}

	return user, nil
}

// Login return user information in case of correct email and password.
func (s *UserService) Login(input LoginInput) (*models.User, error) {
	pass, _ := models.NewPassword(input.Password)

	user := &models.User{
		Email:    input.Email,
		Password: pass,
	}

	v := validator.New()
	if user.Validate(v); !v.Valid() {
		return nil, v.Err()
	}

	user, err := s.repos.Users.FindByEmail(user.Email)
	if err != nil {
		switch {
		// for security sake. if email not found we return unauthorized too.
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrUnauthorized
		default:
			return nil, fmt.Errorf("user login failed: %w", err)
		}
	}

	if !user.Password.ComparePassword(input.Password) {
		return nil, ErrUnauthorized
	}

	return user, nil
}

// FindByID returns user by ID. This is intended for auth middleware
// retrievals. For login use Login method instead
func (s *UserService) FindByID(id int) (*models.User, error) {
	user, err := s.repos.Users.FindByID(id)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, ErrUnauthorized
		default:
			return nil, fmt.Errorf("find user by ID failed: %w", err)
		}
	}

	return user, nil
}

type SignupCustomerInput struct {
	Email    string
	Password string
}

type SignupMerchantInput struct {
	Email     string
	Password  string
	StoreName string
}

type SignupAdminInput struct {
	Admin    *models.User
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}
