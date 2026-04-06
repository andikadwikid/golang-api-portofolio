package utils

import "github.com/go-playground/validator/v10"

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		errors["error"] = err.Error()
		return errors
	}

	for _, e := range validationErrors {
		field := e.Field()

		switch e.Tag() {
		case "required":
			errors[field] = field + " is required"
		case "email":
			errors[field] = "Invalid email format"
		case "min":
			errors[field] = field + " must be at least " + e.Param() + " characters"
		case "max":
			errors[field] = field + " must be at most " + e.Param() + " characters"
		case "unique":
			errors[field] = field + " is already used"
		default:
			errors[field] = field + " is invalid"
		}
	}

	return errors
}
