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
	decode encoding.DecodeFunc
}

// NewRequest creates a new request.
func NewRequest(f map[string]string, r io.Reader, d encoding.DecodeFunc) *Request {
	return &Request{Fields: f, Raw: r, decode: d}
}

// Decode the raw data by using the provided decoder.
func (r *Request) Decode(v interface{}) error {
	return r.decode(r.Raw, v)
}

// Response definition of the sync response model.
type Response struct {
	Payload interface{}
}

// NewResponse creates a new response.
func NewResponse(p interface{}) *Response {
	return &Response{Payload: p}
}

// ProcessorFunc definition of a function type for processing sync requests.
type ProcessorFunc func(context.Context, *Request) (*Response, error)
