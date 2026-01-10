package models

import "github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCustomer Role = "customer"
	RoleMerchant Role = "merchant"
)

func (r Role) Validate(v *validator.Validator) {
	switch r {
	case RoleAdmin, RoleCustomer, RoleMerchant:
	default:
		v.AddError("role", "invalid role")
	}
}
