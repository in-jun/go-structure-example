package validate

import (
	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Struct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fieldErrors := make(map[string]string)
			for _, e := range validationErrors {
				fieldErrors[e.Field()] = getErrorMsg(e)
			}
			return errors.BadRequest("Invalid input").WithDetails(fieldErrors)
		}
		return errors.BadRequest("Invalid input")
	}
	return nil
}

func getErrorMsg(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	default:
		return "Invalid value"
	}
}
