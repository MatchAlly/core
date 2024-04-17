package helpers

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	return &Validator{
		validator: v,
	}
}

func (cv *Validator) Validate(s any) error {
	if err := cv.validator.Struct(s); err != nil {
		return errors.Wrap(err, "validation failed")
	}
	return nil
}
