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
	Raw     io.ReadCloser
	decode  encoding.Decode
}

// NewRequest creates a new request item
func NewRequest(h map[string]string, f map[string]string, r io.ReadCloser, d encoding.Decode) *Request {
	return &Request{
		Headers: h,
		Fields:  f,
		Raw:     r,
		decode:  d,
	}
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

type ValidationError struct {
}

func (e *ValidationError) Error() string {
	return ""
}

type UnauthorizedError struct {
}

func (e *UnauthorizedError) Error() string {
	return ""
}

type ForbiddenError struct {
}

func (e *ForbiddenError) Error() string {
	return ""
}

type NotFoundError struct {
}

func (e *NotFoundError) Error() string {
	return ""
}
