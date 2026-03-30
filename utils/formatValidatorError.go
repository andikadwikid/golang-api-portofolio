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
			errors[field] = field + " wajib diisi"
		case "email":
			errors[field] = "Format email tidak valid"
		case "min":
			errors[field] = field + " minimal " + e.Param() + " karakter"
		default:
			errors[field] = field + " tidak valid"
		}
	}

	return errors
}
