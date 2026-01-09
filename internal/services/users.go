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

// Signup used for creating new users. note that the users can only be
// customers or merchants. for creating admins use SignupAdmin.
func (s *UserService) Signup(input SignupInput) (*models.User, error) {
	pass, _ := models.NewPassword(input.Password)

	user := &models.User{
		Email:    input.Email,
		Role:     input.Role,
		Password: pass,
	}

	// validate that user is not an admin
	v := validator.New()
	if user.Role != "customer" && user.Role != "merchant" {
		v.AddError("role", "role must be (cutomer or merchant)")
		return nil, v.Err()
	}

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

// SignupAdmin creates new users from any type including admins. note that the
// request must be sent from an admin user to be valid.
func (s *UserService) SignupAdmin(input SignupAdminInput) (*models.User, error) {
	if input.Admin.Role != "admin" {
		return nil, ErrUnauthorized
	}

	pass, _ := models.NewPassword(input.Password)

	user := &models.User{
		Email:    input.Email,
		Role:     input.Role,
		Password: pass,
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

type SignupInput struct {
	Email    string
	Password string
	Role     string
}

type SignupAdminInput struct {
	Admin    *models.User
	Email    string
	Password string
	Role     string
}

type LoginInput struct {
	Email    string
	Password string
}
