package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e CustomError) Error() string {
	if e.Details != nil {
		if details, err := json.Marshal(e.Details); err == nil {
			return fmt.Sprintf("status: %d, message: %s, details: %s", e.Status, e.Message, string(details))
		}
	}
	return fmt.Sprintf("status: %d, message: %s", e.Status, e.Message)
}

func (e CustomError) WithDetails(details interface{}) CustomError {
	e.Details = details
	return e
}

func (e CustomError) Is(target error) bool {
	t, ok := target.(CustomError)
	if !ok {
		return false
	}
	return e.Status == t.Status
}

func New(status int, message string) CustomError {
	return CustomError{
		Status:  status,
		Message: message,
	}
}

func FromError(err error) CustomError {
	if customErr, ok := err.(CustomError); ok {
		return customErr
	}
	return Internal("Internal Server Error")
}

func BadRequest(message string) CustomError {
	return CustomError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func Unauthorized(message string) CustomError {
	return CustomError{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func PaymentRequired(message string) CustomError {
	return CustomError{
		Status:  http.StatusPaymentRequired,
		Message: message,
	}
}

func Forbidden(message string) CustomError {
	return CustomError{
		Status:  http.StatusForbidden,
		Message: message,
	}
}

func NotFound(message string) CustomError {
	return CustomError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func MethodNotAllowed(message string) CustomError {
	return CustomError{
		Status:  http.StatusMethodNotAllowed,
		Message: message,
	}
}

func NotAcceptable(message string) CustomError {
	return CustomError{
		Status:  http.StatusNotAcceptable,
		Message: message,
	}
}

func ProxyAuthRequired(message string) CustomError {
	return CustomError{
		Status:  http.StatusProxyAuthRequired,
		Message: message,
	}
}

func RequestTimeout(message string) CustomError {
	return CustomError{
		Status:  http.StatusRequestTimeout,
		Message: message,
	}
}

func Conflict(message string) CustomError {
	return CustomError{
		Status:  http.StatusConflict,
		Message: message,
	}
}

func Gone(message string) CustomError {
	return CustomError{
		Status:  http.StatusGone,
		Message: message,
	}
}

func LengthRequired(message string) CustomError {
	return CustomError{
		Status:  http.StatusLengthRequired,
		Message: message,
	}
}

func PreconditionFailed(message string) CustomError {
	return CustomError{
		Status:  http.StatusPreconditionFailed,
		Message: message,
	}
}

func RequestEntityTooLarge(message string) CustomError {
	return CustomError{
		Status:  http.StatusRequestEntityTooLarge,
		Message: message,
	}
}

func RequestURITooLong(message string) CustomError {
	return CustomError{
		Status:  http.StatusRequestURITooLong,
		Message: message,
	}
}

func UnsupportedMediaType(message string) CustomError {
	return CustomError{
		Status:  http.StatusUnsupportedMediaType,
		Message: message,
	}
}

func RequestedRangeNotSatisfiable(message string) CustomError {
	return CustomError{
		Status:  http.StatusRequestedRangeNotSatisfiable,
		Message: message,
	}
}

func ExpectationFailed(message string) CustomError {
	return CustomError{
		Status:  http.StatusExpectationFailed,
		Message: message,
	}
}

func Teapot(message string) CustomError {
	return CustomError{
		Status:  http.StatusTeapot,
		Message: message,
	}
}

func UnprocessableEntity(message string) CustomError {
	return CustomError{
		Status:  http.StatusUnprocessableEntity,
		Message: message,
	}
}

func Locked(message string) CustomError {
	return CustomError{
		Status:  http.StatusLocked,
		Message: message,
	}
}

func FailedDependency(message string) CustomError {
	return CustomError{
		Status:  http.StatusFailedDependency,
		Message: message,
	}
}

func TooEarly(message string) CustomError {
	return CustomError{
		Status:  http.StatusTooEarly,
		Message: message,
	}
}

func UpgradeRequired(message string) CustomError {
	return CustomError{
		Status:  http.StatusUpgradeRequired,
		Message: message,
	}
}

func PreconditionRequired(message string) CustomError {
	return CustomError{
		Status:  http.StatusPreconditionRequired,
		Message: message,
	}
}

func TooManyRequests(message string) CustomError {
	return CustomError{
		Status:  http.StatusTooManyRequests,
		Message: message,
	}
}

func RequestHeaderFieldsTooLarge(message string) CustomError {
	return CustomError{
		Status:  http.StatusRequestHeaderFieldsTooLarge,
		Message: message,
	}
}

func UnavailableForLegalReasons(message string) CustomError {
	return CustomError{
		Status:  http.StatusUnavailableForLegalReasons,
		Message: message,
	}
}

func Internal(message string) CustomError {
	return CustomError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

func NotImplemented(message string) CustomError {
	return CustomError{
		Status:  http.StatusNotImplemented,
		Message: message,
	}
}

func BadGateway(message string) CustomError {
	return CustomError{
		Status:  http.StatusBadGateway,
		Message: message,
	}
}

func ServiceUnavailable(message string) CustomError {
	return CustomError{
		Status:  http.StatusServiceUnavailable,
		Message: message,
	}
}

func GatewayTimeout(message string) CustomError {
	return CustomError{
		Status:  http.StatusGatewayTimeout,
		Message: message,
	}
}

func HTTPVersionNotSupported(message string) CustomError {
	return CustomError{
		Status:  http.StatusHTTPVersionNotSupported,
		Message: message,
	}
}

func VariantAlsoNegotiates(message string) CustomError {
	return CustomError{
		Status:  http.StatusVariantAlsoNegotiates,
		Message: message,
	}
}

func InsufficientStorage(message string) CustomError {
	return CustomError{
		Status:  http.StatusInsufficientStorage,
		Message: message,
	}
}

func LoopDetected(message string) CustomError {
	return CustomError{
		Status:  http.StatusLoopDetected,
		Message: message,
	}
}

func NotExtended(message string) CustomError {
	return CustomError{
		Status:  http.StatusNotExtended,
		Message: message,
	}
}

func NetworkAuthenticationRequired(message string) CustomError {
	return CustomError{
		Status:  http.StatusNetworkAuthenticationRequired,
		Message: message,
	}
}
