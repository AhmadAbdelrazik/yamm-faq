package models

import (
	"strings"

	"github.com/AhmadAbdelrazik/yamm_faq/pkg/validator"
)

type Translation struct {
	ID       int
	FAQID    int
	Language Language
	Question string
	Answer   string
}

func (t *Translation) Validate(v *validator.Validator) {
	v.Check(t.ID >= 0, "id", "invalid id")

	t.Language.Validate(v)

	v.Check(len(strings.TrimSpace(t.Question)) > 0, "question", "required")
	v.Check(len(t.Question) <= 1000, "question", "must not be longer than 1000 characters")

	v.Check(len(strings.TrimSpace(t.Answer)) > 0, "answer", "required")
	v.Check(len(t.Answer) <= 3000, "answer", "must not be longer than 3000 characters")
}
