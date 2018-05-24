package sync

import (
	"context"
	"io"
)

// Unmarshaller definition of a function for unmarshalling a model.
type Unmarshaller func(data io.ReadCloser, v interface{}) error

// Request definition of the sync request model.
type Request struct {
	Headers      map[string]string
	Fields       map[string]string
	Raw          io.ReadCloser
	unmarshaller Unmarshaller
}

// NewRequest creates a new request item
func NewRequest(h map[string]string, f map[string]string, r io.ReadCloser, u Unmarshaller) *Request {
	return &Request{
		Headers:      h,
		Fields:       f,
		Raw:          r,
		unmarshaller: u,
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
