package jsonvalidator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type JSONValidator struct {
	*validator.Validate
}

func New(v *validator.Validate) *JSONValidator {
	return &JSONValidator{
		v,
	}
}

func (v *JSONValidator) Struct(s interface{}) error {
	errs := v.Validate.Struct(s)

	if errs == nil {
		return nil
	}

	newErrMes := ""
	for _, v := range errs.(validator.ValidationErrors) {
		fmt.Print(v.Type())
		if v.Tag() == "required" {
			newErrMes += fmt.Sprintf("Field %s must be provided\n", v.Field())
		}
		if v.Tag() == "email" {
			newErrMes += fmt.Sprintf("Field %s must contains email\n", v.Field())
		}
		if v.Tag() == "min" {
			newErrMes += fmt.Sprintf("Minimal lenght for field %s is %v\n", v.Field(), v.Param())
		}
		if v.Tag() == "max" {
			newErrMes += fmt.Sprintf("Maximum lenght for field %s is %v\n", v.Field(), v.Param())
		}
		fmt.Print(v.Tag())
	}

	return errors.New(newErrMes)
}
