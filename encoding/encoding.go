package encoding

import "io"

// Decode definition of a JSON decoding function.
type Decode func(data io.Reader, v interface{}) error

// Encode definition of a JSON encoding function.
type Encode func(v interface{}) ([]byte, error)
