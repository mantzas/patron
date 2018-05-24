package json

import (
	"encoding/json"
	"io"
)

// Decode a JSON input in the form of a read closer.
func Decode(data io.ReadCloser, v interface{}) error {
	defer data.Close()
	return json.NewDecoder(data).Decode(v)
}

// Encode a model to JSON.
func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
