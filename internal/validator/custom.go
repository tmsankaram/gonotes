package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func TitleNoTest(fl validator.FieldLevel) bool {
	title := fl.Field().String()
	return !strings.Contains(strings.ToLower(title), "test")
}
