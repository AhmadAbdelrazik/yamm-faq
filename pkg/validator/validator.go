package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	EmailRX    = regexp.MustCompile(`^([a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6})*$`)
	LowerRX    = regexp.MustCompile(`[a-z]`)
	AlphanumRX = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	UpperRX    = regexp.MustCompile(`[A-Z]`)
	NumberRX   = regexp.MustCompile(`[0-9]`)
	SpecialRX  = regexp.MustCompile(`[!@#$%^&*]`)
)

// Validator used for producing errors in key-value format. it's mainly used
// for struct validation.
type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) AddError(key, value string) {
	if _, ok := v.Errors[key]; !ok {
		v.Errors[key] = value
	}
}

// Check if condition is true. If false add an error
func (v *Validator) Check(condition bool, key, value string) {
	if !condition {
		v.AddError(key, value)
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) Err() error {
	if v.Valid() {
		return nil
	}

	var errorString strings.Builder
	for k, v := range v.Errors {
		fmt.Fprintf(&errorString, "%q:%q ", k, v)
	}

	return errors.New(errorString.String())
}
