package encoding

import (
	"io"
)

const (
	// ContentTypeHeader for defining content types
	ContentTypeHeader string = "Content-Type"
)

// Decode definition of a JSON decoding function from a reader.
type Decode func(data io.Reader, v interface{}) error

// DecodeRaw definition of a JSON decoding function from a byte slice.
type DecodeRaw func(data []byte, v interface{}) error

// Encode definition of a JSON encoding function.
type Encode func(v interface{}) ([]byte, error)
