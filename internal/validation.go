package utils

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func ParseValidationErrors(err error) string {
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, err.Field()+" is "+err.Tag())
	}
	return strings.Join(errors, ", ")
}
