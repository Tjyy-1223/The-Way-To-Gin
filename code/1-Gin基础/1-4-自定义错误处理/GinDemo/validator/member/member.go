package member

import "github.com/go-playground/validator/v10"

func NameValid(f1 validator.FieldLevel) bool {
	s := f1.Field().String()
	if s == "admin" {
		return false
	}
	return true
}
