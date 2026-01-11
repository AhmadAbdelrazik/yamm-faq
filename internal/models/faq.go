package models

import (
	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

type FAQ struct {
	ID              int
	Category        FAQCategory
	StoreID         int
	IsGlobal        bool
	DefaultLanguage Language
	Translations    []Translation
}

func (f *FAQ) Validate(v *validator.Validator) {
	v.Check(f.ID >= 0, "id", "invalid id")

	f.Category.Validate(v)
	for _, t := range f.Translations {
		t.Validate(v)
		if !v.Valid() {
			break
		}
	}

	v.Check(f.StoreID >= 0, "store_id", "invalid id")
}
