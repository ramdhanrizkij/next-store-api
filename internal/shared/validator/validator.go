package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError formats validator.ValidationErrors into a user-friendly map
func FormatValidationError(err error) map[string]string {
	res := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			res[field] = msgForTag(fe)
		}
		return res
	}

	res["error"] = err.Error()
	return res
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters long", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Failed validation on tag '%s'", fe.Tag())
	}
}
