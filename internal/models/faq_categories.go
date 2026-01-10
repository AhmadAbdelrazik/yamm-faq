package models

import (
	"strings"

	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

type FAQCategory struct {
	ID   int
	Name string
}

func (f *FAQCategory) Validate(v *validator.Validator) {
	v.Check(f.ID >= 0, "id", "invalid id")

	v.Check(len(strings.TrimSpace(f.Name)) > 0, "category_name", "required")
	v.Check(len(f.Name) <= 20, "category_name", "must not be more than 20 bytes")
}
