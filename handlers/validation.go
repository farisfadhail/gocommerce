package handlers

import "github.com/go-playground/validator/v10"

func FormatValidationError(err error) []string {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, err.Field()+" is "+err.Tag())
	}

	return errors
}
