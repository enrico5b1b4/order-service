package app

import (
	"fmt"
	"net/http"

	"github.com/enrico5b1b4/order-service/errors"
	"github.com/enrico5b1b4/order-service/order"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		fieldValidationErrors := make([]errors.FieldValidatonError, len(validationErr))
		for i := range validationErr {
			fieldValidationErrors[i].Key = validationErr[i].Field()
			fieldValidationErrors[i].Reason = fmt.Sprintf("Field contains errors: %s", validationErr[i].Tag())
		}

		validationError := errors.NewValidationError("Validation failed", fieldValidationErrors)

		c.JSON(http.StatusBadRequest, validationError) // #nosec
		return
	}

	if orderServiceError, ok := err.(*errors.Error); ok {
		statusCode := ErrorToStatusCode(orderServiceError.Code)

		c.JSON(statusCode, orderServiceError) // #nosec
		return
	}

	c.JSON(http.StatusInternalServerError, errors.GeneralError) // #nosec
}

var ErrorToStatusCodeMap = map[string]int{
	errors.GeneralError.Code:                   http.StatusInternalServerError,
	order.OrderNotFoundError.Code:              http.StatusNotFound,
	order.OrderAlreadyExistsError.Code:         http.StatusBadRequest,
	order.OrderAlreadyBeingProcessedError.Code: http.StatusBadRequest,
	order.OrderAlreadyCompleteError.Code:       http.StatusBadRequest,
}

func ErrorToStatusCode(code string) int {
	if statusCode, ok := ErrorToStatusCodeMap[code]; ok {
		return statusCode
	}
	return http.StatusInternalServerError
}
