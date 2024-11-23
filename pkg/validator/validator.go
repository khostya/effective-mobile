package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	v         *validator.Validate
	passwdErr error
}

func NewValidate() (*Validator, error) {
	v := Validator{v: validator.New(validator.WithRequiredStructEnabled())}

	return &v, nil
}

func (v *Validator) Struct(i any) error {
	return v.v.Struct(i)
}
