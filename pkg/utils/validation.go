package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError converts raw Go struct validator errors into user-friendly strings.
func FormatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var out []string
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				out = append(out, field+" must be filled")
			case "email":
				out = append(out, field+" must be a valid email address")
			case "min":
				out = append(out, field+" must be at least "+fe.Param()+" characters")
			case "max":
				out = append(out, field+" must be at most "+fe.Param()+" characters")
			case "oneof":
				out = append(out, field+" must be one of: "+fe.Param())
			default:
				out = append(out, field+" is invalid")
			}
		}
		return strings.Join(out, ", ")
	}
	return "invalid request parameters"
}
