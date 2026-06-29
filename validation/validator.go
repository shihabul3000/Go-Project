package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shihabul3000/Go-Project/apperrors"
)

type CustomValidator struct {
	validate *validator.Validate
}

func New() *CustomValidator {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validate: validate}
}

func (cv *CustomValidator) Validate(value interface{}) error {
	if err := cv.validate.Struct(value); err != nil {
		fields := map[string]string{}
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				fields[fieldErr.Field()] = messageFor(fieldErr)
			}
			return &apperrors.ValidationError{Fields: fields}
		}
		return &apperrors.ValidationError{Fields: map[string]string{"request": err.Error()}}
	}

	return nil
}

func messageFor(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", err.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", err.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", err.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", err.Param())
	default:
		return "is invalid"
	}
}
