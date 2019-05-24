package request

import (
	"reflect"
	"strings"

	"github.com/enrico5b1b4/order-service/order/completeorder"
	"gopkg.in/go-playground/validator.v9"
)

type Validator struct {
	Validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}

func NewRequestValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	_ = validate.RegisterValidation("completeOrderStatusValidation", completeorder.CompleteOrderStatusValidation)

	return validate
}
