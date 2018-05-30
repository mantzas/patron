package sync

import (
	"context"
	"io"

	"github.com/mantzas/patron/encoding"
)

// Request definition of the sync request model.
type Request struct {
	Headers map[string]string
	Fields  map[string]string
	Raw     io.Reader
	decode  encoding.Decode
}

// NewRequest creates a new request item
func NewRequest(h map[string]string, f map[string]string, r io.Reader, d encoding.Decode) *Request {
	return &Request{h, f, r, d}
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

// Handler definition of a generic sync handler.
type Handler interface {
	Handle(context.Context, *Request) (*Response, error)
}

// ValidationError defines a validation error.
type ValidationError struct {
	err string
}

func (e *ValidationError) Error() string {
	return e.err
}

// UnauthorizedError defines a authorization error.
type UnauthorizedError struct {
	err string
}

func (e *UnauthorizedError) Error() string {
	return e.err
}

// ForbiddenError defines a access error.
type ForbiddenError struct {
	err string
}

func (e *ForbiddenError) Error() string {
	return e.err
}

// NotFoundError defines a not found error.
type NotFoundError struct {
	err string
}

func (e *NotFoundError) Error() string {
	return e.err
}

// ServiceUnavailableError defines a service unavailable error.
type ServiceUnavailableError struct {
	err string
}

func (e *ServiceUnavailableError) Error() string {
	return e.err
}
