package errors

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e CustomError) Error() string {
	return fmt.Sprintf("status: %d, code: %s, message: %s", e.Status, e.Code, e.Message)
}

func (e CustomError) Is(target error) bool {
	t, ok := target.(CustomError)
	if !ok {
		return false
	}
	return e.Status == t.Status
}

func BadRequest(message string) CustomError {
	return CustomError{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: message}
}

func ValidationError(message string) CustomError {
	return CustomError{Status: http.StatusBadRequest, Code: "VALIDATION_ERROR", Message: message}
}

func Unauthorized(message string) CustomError {
	return CustomError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: message}
}

func Forbidden(message string) CustomError {
	return CustomError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: message}
}

func NotFound(message string) CustomError {
	return CustomError{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: message}
}

func Conflict(message string) CustomError {
	return CustomError{Status: http.StatusConflict, Code: "CONFLICT", Message: message}
}

func TooManyRequests(message string) CustomError {
	return CustomError{Status: http.StatusTooManyRequests, Code: "RATE_LIMIT_EXCEEDED", Message: message}
}

func Internal(message string) CustomError {
	return CustomError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: message}
}
