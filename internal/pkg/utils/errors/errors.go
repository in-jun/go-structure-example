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

func New(status int, code, message string) CustomError {
	return CustomError{Status: status, Code: code, Message: message}
}

func BadRequest(message string) CustomError {
	return CustomError{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: message}
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

func UnprocessableEntity(message string) CustomError {
	return CustomError{Status: http.StatusUnprocessableEntity, Code: "UNPROCESSABLE_ENTITY", Message: message}
}

func TooManyRequests(message string) CustomError {
	return CustomError{Status: http.StatusTooManyRequests, Code: "TOO_MANY_REQUESTS", Message: message}
}

func Internal(message string) CustomError {
	return CustomError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: message}
}

func NotImplemented(message string) CustomError {
	return CustomError{Status: http.StatusNotImplemented, Code: "NOT_IMPLEMENTED", Message: message}
}

func ServiceUnavailable(message string) CustomError {
	return CustomError{Status: http.StatusServiceUnavailable, Code: "SERVICE_UNAVAILABLE", Message: message}
}
