package custom

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Custom email validator
func EmailValidator(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}

// Initialize validator and register custom validators
func InitValidator() (*validator.Validate, error) {
	validate := validator.New()
	if err := validate.RegisterValidation("emailValidator", EmailValidator); err != nil {
		return nil, err
	}
	return validate, nil
}
