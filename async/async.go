package async

import (
	"context"

	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/pkg/errors"
)

// Processor interface for implementing processing of messages
type Processor interface {
	Process(context.Context, *Message) error
}

// Message definition of a async message.
type Message struct {
	data   []byte
	decode encoding.DecodeRaw
}

// NewMessage creates a new message.
func NewMessage(d []byte, dec encoding.DecodeRaw) *Message {
	return &Message{d, dec}
}

// Decode a the raw message into the given value.
func (m *Message) Decode(v interface{}) error {
	return m.decode(m.data, v)
}

// DetermineDecoder determines the decoder based on the content type.
func DetermineDecoder(contentType string) (encoding.DecodeRaw, error) {
	switch contentType {
	case json.ContentType, json.ContentTypeCharset:
		return json.DecodeRaw, nil
	}
	return nil, errors.Errorf("accept header %s is unsupported", contentType)
}
