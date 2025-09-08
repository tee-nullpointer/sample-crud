package customerrors

import (
	"fmt"
	"net/http"
)

const (
	CodeBadRequest = "40"
	CodeNotFound   = "44"
)

type CustomError struct {
	Code       string
	HttpStatus int
	Message    string
	Details    string
}

func (e CustomError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

func NewCustomError(httpStatus int, code string, message string, detail string) *CustomError {
	return &CustomError{
		HttpStatus: httpStatus,
		Code:       code,
		Message:    message,
		Details:    detail,
	}
}

func NewBadRequestError(message string, detail string) *CustomError {
	return NewCustomError(http.StatusBadRequest, CodeBadRequest, message, detail)
}

func NewNotFoundError(message string, detail string) *CustomError {
	return NewCustomError(http.StatusNotFound, CodeNotFound, message, detail)
}
