package validator

import (
	"reflect"

	pkgErr "expensez/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	Validator *CustomValidator
)

type CustomValidator struct {
	validator *validator.Validate
	errorMap  map[string]error
}

func Load(errorMap map[string]error) error {
	Validator = &CustomValidator{validator: validator.New(), errorMap: errorMap}
	return nil
}

func (cv *CustomValidator) Validate(i any) []error {
	var err error

	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		err = cv.validator.Var(i, "required,min=1,dive")
	default:
		err = cv.validator.Struct(i) // reflect.Struct:

	}
	if err != nil {
		var errs []error
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			if er, ok := cv.errorMap[err.Field()+"_"+err.Tag()]; ok {
				errs = append(errs, er)
			} else {
				errs = append(errs, pkgErr.New(pkgErr.VALIDATION_ERROR, err.Error()))
			}
		}
		return errs
	}
	return nil
}

func Validate(c echo.Context, modal any) []error {
	if err := c.Bind(modal); err != nil {
		return []error{pkgErr.New(pkgErr.BAD_REQUEST, err.Error())}
	}
	val := &CustomValidator{validator: validator.New()}
	return val.Validate(modal)
}

/*
Request body:
err := (&DefaultBinder{}).BindBody(c, &payload)

Query parameters:
err := (&DefaultBinder{}).BindQueryParams(c, &payload)

Path parameters:
err := (&DefaultBinder{}).BindPathParams(c, &payload)

Header parameters:
err := (&DefaultBinder{}).BindHeaders(c, &payload)
*/
