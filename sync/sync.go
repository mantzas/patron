package sync

import (
	"context"
	"io"

	"github.com/mantzas/patron/encoding"
)

// Request definition of the sync request model.
type Request struct {
	Fields map[string]string
	Raw    io.Reader
	decode encoding.Decode
}

// NewRequest creates a new request item
func NewRequest(f map[string]string, r io.Reader, d encoding.Decode) *Request {
	return &Request{f, r, d}
}

// Decode a the raw message into the given value.
func (r *Request) Decode(v interface{}) error {
	return r.decode(r.Raw, v)
}

// Response definition of the sync response model.
type Response struct {
	Payload interface{}
}

// NewResponse creates a new response.
func NewResponse(p interface{}) *Response {
	return &Response{p}
}

// ProcessorFunc defines a function type for processing sync requests.
type ProcessorFunc func(context.Context, *Request) (*Response, error)

// ValidationError defines a validation error.
type ValidationError struct {
	err string
}

func (e *ValidationError) Error() string {
	return e.err
}

// NewValidationError creates a new validation error.
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{msg}
}

// UnauthorizedError defines a authorization error.
type UnauthorizedError struct {
	err string
}

func (e *UnauthorizedError) Error() string {
	return e.err
}

// NewUnauthorizedError creates a new unauthorized error.
func NewUnauthorizedError(msg string) *UnauthorizedError {
	return &UnauthorizedError{msg}
}

// ForbiddenError defines a access error.
type ForbiddenError struct {
	err string
}

func (e *ForbiddenError) Error() string {
	return e.err
}

// NewForbiddenError creates a new forbidden error.
func NewForbiddenError(msg string) *ForbiddenError {
	return &ForbiddenError{msg}
}

// NotFoundError defines a not found error.
type NotFoundError struct {
	err string
}

func (e *NotFoundError) Error() string {
	return e.err
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{msg}
}

// ServiceUnavailableError defines a service unavailable error.
type ServiceUnavailableError struct {
	err string
}

func (e *ServiceUnavailableError) Error() string {
	return e.err
}

// NewServiceUnavailableError creates a new service unavailable error.
func NewServiceUnavailableError(msg string) *ServiceUnavailableError {
	return &ServiceUnavailableError{msg}
}
