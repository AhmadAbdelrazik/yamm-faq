package models

import (
	"strings"

	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    string
	Role     Role
	Password *Password
}

func (u *User) Validate(v *validator.Validator) {
	v.Check(u.ID >= 0, "id", "invalid id")

	v.Check(len(strings.TrimSpace(u.Email)) > 0, "email", "required")
	v.Check(validator.EmailRX.MatchString(u.Email), "email", "invalid email form")

	u.Role.Validate(v)
	u.Password.Validate(v)
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsMerchant() bool {
	return u.Role == RoleMerchant
}

func (u *User) IsCustomer() bool {
	return u.Role == RoleCustomer
}

type Password struct {
	password *string
	Hash     []byte
}

func NewPassword(password string) (*Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Password{
		password: &password,
		Hash:     hash,
	}, nil
}

func (i *Password) Validate(v *validator.Validator) {

	if i.password == nil {
		v.AddError("password", "required")
		return
	}
	password := *i.password

	v.Check(len(strings.TrimSpace(password)) > 0, "password", "required")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters")
	v.Check(len(password) <= 50, "password", "must be at most 50 characters")
	v.Check(
		validator.LowerRX.MatchString(password),
		"password",
		"must contain at least 1 lowercase character",
	)
	v.Check(
		validator.UpperRX.MatchString(password),
		"password",
		"must contain at least 1 uppercase character",
	)
	v.Check(
		validator.NumberRX.MatchString(password),
		"password",
		"must contain at least a number",
	)
	v.Check(
		validator.SpecialRX.MatchString(password),
		"password",
		"must contain at least 1 special character ( !@#$%&* )",
	)
}

func (p *Password) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(password)) == nil
}
