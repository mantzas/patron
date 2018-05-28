package json

import (
	"encoding/json"
	"io"
)

// Decode a JSON input in the form of a read closer.
func Decode(data io.Reader, v interface{}) error {
	return json.NewDecoder(data).Decode(v)
}

// Encode a model to JSON.
func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
