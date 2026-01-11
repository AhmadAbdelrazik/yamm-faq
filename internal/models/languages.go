package models

import "github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"

type Language string

const (
	LanguageArabic  Language = "ar"
	LanguageEnglish Language = "en"
	LanguageSpanish Language = "es"
	LanguageGerman  Language = "de"
)

func (l Language) Validate(v *validator.Validator) {
	switch l {
	case LanguageArabic, LanguageEnglish, LanguageSpanish, LanguageGerman:
	default:
		v.AddError("language", "invalid language")
	}
}
