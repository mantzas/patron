package http

import (
	gohttp "net/http"
)

// Error defines an abstract struct that can represent several types of HTTP errors.
type Error struct {
	code    int
	payload interface{}
}

// Error returns the actual message of the error.
func (e *Error) Error() string {
	err, _ := e.payload.(string)
	return err
}

// NewValidationError creates a new validation error with default payload.
func NewValidationError() *Error {
	return &Error{gohttp.StatusBadRequest, gohttp.StatusText(gohttp.StatusBadRequest)}
}

// NewValidationErrorWithPayload creates a new validation error with the specified payload.
func NewValidationErrorWithPayload(payload interface{}) *Error {
	return &Error{gohttp.StatusBadRequest, payload}
}

// NewUnauthorizedError creates a new validation error with default payload.
func NewUnauthorizedError() *Error {
	return &Error{gohttp.StatusUnauthorized, gohttp.StatusText(gohttp.StatusUnauthorized)}
}

// NewUnauthorizedErrorWithPayload creates a new unauthorized error with the specified payload.
func NewUnauthorizedErrorWithPayload(payload interface{}) *Error {
	return &Error{gohttp.StatusUnauthorized, payload}
}

// NewForbiddenError creates a new forbidden error with default payload.
func NewForbiddenError() *Error {
	return &Error{gohttp.StatusForbidden, gohttp.StatusText(gohttp.StatusForbidden)}
}

// NewForbiddenErrorWithPayload creates a new forbidden error with the specified payload.
func NewForbiddenErrorWithPayload(payload interface{}) *Error {
	return &Error{gohttp.StatusForbidden, payload}
}

// NewNotFoundError creates a new not found error with default payload.
func NewNotFoundError() *Error {
	return &Error{gohttp.StatusNotFound, gohttp.StatusText(gohttp.StatusNotFound)}
}

// NewNotFoundErrorWithPayload creates a new not found error with the specified payload.
func NewNotFoundErrorWithPayload(payload interface{}) *Error {
	return &Error{gohttp.StatusNotFound, payload}
}

// NewServiceUnavailableError creates a new service unavailable error with default payload.
func NewServiceUnavailableError() *Error {
	return &Error{gohttp.StatusServiceUnavailable, gohttp.StatusText(gohttp.StatusServiceUnavailable)}
}

// NewServiceUnavailableErrorWithPayload creates a new service unavailable error with the specified payload.
func NewServiceUnavailableErrorWithPayload(payload interface{}) *Error {
	return &Error{gohttp.StatusServiceUnavailable, payload}
}

// NewError creates a new error with default Internal Server Error payload.
func NewError() *Error {
	return &Error{gohttp.StatusInternalServerError, gohttp.StatusText(gohttp.StatusInternalServerError)}
}

// NewErrorWithCodeAndPayload creates a fully customizable error with the specified status code and payload.
func NewErrorWithCodeAndPayload(code int, payload interface{}) *Error {
	return &Error{code, payload}
}
