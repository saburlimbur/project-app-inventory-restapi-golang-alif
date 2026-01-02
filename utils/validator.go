package utils

import (
	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateErrors(data any) ([]FieldError, error) {
	validate := validator.New()

	err := validate.Struct(data)
	if err == nil {
		return nil, nil
	}

	var errors []FieldError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			var message string
			switch err.Tag() {
			case "required":
				message = err.Field() + " is required"
			case "min":
				message = err.Field() + " must be at least " + err.Param() + " characters"
			case "max":
				message = err.Field() + " must be at most " + err.Param() + " characters"
			case "email":
				message = err.Field() + " must be a valid email address"
			case "alphanum":
				message = err.Field() + " must contain only alphanumeric characters"
			case "oneof":
				message = err.Field() + " must be one of: " + err.Param()
			default:
				message = err.Field() + " is invalid"
			}

			errors = append(errors, FieldError{
				Field:   err.Field(),
				Message: message,
			})
		}
		return errors, err
	}

	// Fallback: return original error if not a validation error
	return nil, err
}
