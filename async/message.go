package async

import (
	"io"

	"github.com/mantzas/patron/encoding"
)

// Message definition of a async message.
type Message struct {
	Headers map[string]string
	Raw     io.Reader
	decode  encoding.Decode
}

// NewMessage creates a new message.
func NewMessage(h map[string]string, r io.Reader, d encoding.Decode) *Message {
	return &Message{h, r, d}
}

// Decode a the raw message into the given value.
func (m *Message) Decode(v interface{}) error {
	return m.decode(m.Raw, v)
}
