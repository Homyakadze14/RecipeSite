package common

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func GetErrMessages(errs error) error {
	if errs == nil {
		return nil
	}

	newErrMes := ""
	var ve validator.ValidationErrors
	if errors.As(errs, &ve) {
		for _, v := range ve {
			if v.Tag() == "required" {
				newErrMes += fmt.Sprintf("Field %s must be provided;", v.Field())
			}
			if v.Tag() == "email" {
				newErrMes += fmt.Sprintf("Field %s must contains email;", v.Field())
			}
			if v.Tag() == "min" {
				newErrMes += fmt.Sprintf("Minimal lenght for field %s is %v;", v.Field(), v.Param())
			}
			if v.Tag() == "max" {
				newErrMes += fmt.Sprintf("Maximum lenght for field %s is %v;", v.Field(), v.Param())
			}
		}
	} else {
		newErrMes = errs.Error()
	}

	return errors.New(newErrMes)
}
