package json

import (
	"encoding/json"
	"io"
)

const (
	// ContentType JSON definition
	ContentType string = "application/json"
	// ContentTypeCharset JSON definition with charset
	ContentTypeCharset string = "application/json; charset=utf-8"
)

// Decode a JSON input in the form of a read.
func Decode(data io.Reader, v interface{}) error {
	return json.NewDecoder(data).Decode(v)
}

// DecodeRaw a JSON input in the form of a byte slice.
func DecodeRaw(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Encode a model to JSON.
func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
