package models

import (
	"strings"

	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

type Store struct {
	ID         int
	MerchantID int
	Name       string
}

func (s *Store) Validate(v *validator.Validator) {
	v.Check(s.ID >= 0, "id", "invalid id")
	v.Check(s.MerchantID >= 0, "merchant_id", "invalid id")

	v.Check(len(strings.TrimSpace(s.Name)) > 0, "name", "required")
	v.Check(len(s.Name) >= 3, "name", "must be at least 3 characters")
	v.Check(len(s.Name) <= 30, "name", "must be at most 30 characters")
}
