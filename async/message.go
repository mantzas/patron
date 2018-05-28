package async

import (
	"github.com/mantzas/patron/encoding"
)

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
