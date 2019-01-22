package http

import (
	gohttp "net/http"
)

// CustomError defines an abstract struct that can represent several types of HTTP errors.
type CustomError struct {
	err     string
	payload interface{}
	code    int
}

// Error returns the actual message of the error.
func (e *CustomError) Error() string {
	return e.err
}

// Payload returns the error payload, which cane be used as HTTP response content.
func (e *CustomError) Payload() interface{} {
	return e.payload
}

// Code returns the status code that corresponds to the specific error type.
func (e *CustomError) Code() int {
	return e.code
}

// NewValidationError creates a new validation error.
func NewValidationError(msg string, payload interface{}) *CustomError {
	return &CustomError{err: msg, code: gohttp.StatusBadRequest, payload: getActualPayload(gohttp.StatusBadRequest, payload)}
}

// NewUnauthorizedError creates a new unauthorized error.
func NewUnauthorizedError(msg string, payload interface{}) *CustomError {
	return &CustomError{err: msg, code: gohttp.StatusUnauthorized, payload: getActualPayload(gohttp.StatusUnauthorized, payload)}
}

// NewForbiddenError creates a new forbidden error.
func NewForbiddenError(msg string, payload interface{}) *CustomError {
	return &CustomError{err: msg, code: gohttp.StatusForbidden, payload: getActualPayload(gohttp.StatusForbidden, payload)}
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(msg string, payload interface{}) *CustomError {
	return &CustomError{err: msg, code: gohttp.StatusNotFound, payload: getActualPayload(gohttp.StatusNotFound, payload)}
}

// NewServiceUnavailableError creates a new service unavailable error.
func NewServiceUnavailableError(msg string, payload interface{}) *CustomError {
	return &CustomError{err: msg, code: gohttp.StatusServiceUnavailable, payload: getActualPayload(gohttp.StatusServiceUnavailable, payload)}
}

func getActualPayload(code int, payload interface{}) interface{} {
	pl := payload
	if pl == nil {
		pl = gohttp.StatusText(code)
	}
	return pl
}
