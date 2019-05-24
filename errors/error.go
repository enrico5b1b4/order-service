package errors

import "fmt"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ValidationError struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Fields  []FieldValidatonError `json:"fields"`
}

type FieldValidatonError struct {
	Key    string `json:"key"`
	Reason string `json:"reason"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s:%s", err.Code, err.Message)
}

func New(message string) *Error {
	return &Error{Code: GeneralError.Code, Message: message}
}

func NewValidationError(message string, fields []FieldValidatonError) *ValidationError {
	return &ValidationError{
		Code:    RequestValidationError.Code,
		Message: message,
		Fields:  fields,
	}
}

var GeneralError = &Error{
	Code:    "0",
	Message: "Something went wrong.",
}

var RequestValidationError = &Error{
	Code: "1",
}
