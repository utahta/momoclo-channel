package validator

import "github.com/go-playground/validator"

var validate = validator.New()

// Validate validates given struct using go-playground/validator
func Validate(src interface{}) error {
	return validate.Struct(src)
}
