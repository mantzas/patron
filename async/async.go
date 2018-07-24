package async

import (
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/pkg/errors"
)

// ProcessorFunc definition of a async processor.
type ProcessorFunc func(MessageI) error

// Message definition of a async message.
type Message struct {
	Data    []byte
	Headers map[string]string
	decode  encoding.DecodeRawFunc
}

// NewMessage creates a new message.
func NewMessage(d []byte, dec encoding.DecodeRawFunc) *Message {
	return &Message{Data: d, decode: dec}
}

// Decode the raw data.
func (m *Message) Decode(v interface{}) error {
	return m.decode(m.Data, v)
}

// DetermineDecoder determines the decoder based on the content type.
func DetermineDecoder(contentType string) (encoding.DecodeRawFunc, error) {
	switch contentType {
	case json.ContentType, json.ContentTypeCharset:
		return json.DecodeRaw, nil
	}
	return nil, errors.Errorf("accept header %s is unsupported", contentType)
}
