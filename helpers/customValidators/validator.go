package customValidators

import "github.com/go-playground/validator"

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	return &CustomValidator{validator: v}
}
