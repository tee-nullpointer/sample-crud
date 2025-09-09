package customerrors

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	CodeBadRequest = "40"
	CodeNotFound   = "44"
)

type CustomError struct {
	Code       string
	HttpStatus int
	GrpcCode   codes.Code
	Message    string
	Details    string
}

func (e CustomError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

func NewCustomError(httpStatus int, grpcCode codes.Code, code string, message string, detail string) *CustomError {
	return &CustomError{
		HttpStatus: httpStatus,
		GrpcCode:   grpcCode,
		Code:       code,
		Message:    message,
		Details:    detail,
	}
}

func NewBadRequestError(message string, detail string) *CustomError {
	return NewCustomError(http.StatusBadRequest, codes.InvalidArgument, CodeBadRequest, message, detail)
}

func NewNotFoundError(message string, detail string) *CustomError {
	return NewCustomError(http.StatusNotFound, codes.NotFound, CodeNotFound, message, detail)
}
