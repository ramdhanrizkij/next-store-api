package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(statusCode int, code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

func Wrap(err error, statusCode int, code, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

var (
	ErrNotFound       = New(http.StatusNotFound, "NOT_FOUND", "Resource not found")
	ErrUnauthorized   = New(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized access")
	ErrForbidden      = New(http.StatusForbidden, "FORBIDDEN", "Forbidden action")
	ErrBadRequest     = New(http.StatusBadRequest, "BAD_REQUEST", "Bad request")
	ErrInternalServer = New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
	ErrConflict       = New(http.StatusConflict, "CONFLICT", "Resource already exists")
	ErrInvalidToken   = New(http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
)

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
