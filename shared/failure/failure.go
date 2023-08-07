package failure

import (
	"fmt"
	"net/http"
)

// Failure is a wrapper for error messages and codes using standard HTTP response codes.
type Failure struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error returns the error code and message in a formatted string.
func (e *Failure) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.Code), e.Message)
}

// BadRequest returns a new Failure with code for bad requests.
func BadRequest(err error) error {
	if err != nil {
		return &Failure{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	return nil
}

// BadRequestFromString returns a new Failure with code for bad requests with message set from string.
func BadRequestFromString(msg string) error {
	return &Failure{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

// Unauthorized returns a new Failure with code for unauthorized requests.
func Unauthorized(msg string) error {
	return &Failure{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}

// InternalError returns a new Failure with code for internal error and message derived from an error interface.
func InternalError(err error) error {
	if err != nil {
		return &Failure{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return nil
}

// Unimplemented returns a new Failure with code for unimplemented method.
func Unimplemented(methodName string) error {
	return &Failure{
		Code:    http.StatusNotImplemented,
		Message: methodName,
	}
}

// NotFound returns a new Failure with code for entity not found.
func NotFound(entityName string) error {
	return &Failure{
		Code:    http.StatusNotFound,
		Message: entityName,
	}
}

// Conflict returns a new Failure with code for conflict situations.
func Conflict(operationName string, entityName string, message string) error {
	return &Failure{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf("%s on %s: %s", operationName, entityName, message),
	}
}

// GetCode returns the error code of an error interface.
func GetCode(err error) int {
	if f, ok := err.(*Failure); ok {
		return f.Code
	}
	return http.StatusInternalServerError
}
