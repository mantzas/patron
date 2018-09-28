package async

import (
	"context"

	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/errors"
)

// ProcessorFunc definition of a async processor.
type ProcessorFunc func(Message) error

// Message interface for defining messages that are handled by the async component.
type Message interface {
	Context() context.Context
	Decode(v interface{}) error
	Ack() error
	Nack() error
}

// ConsumerFactory interface for creating consumers.
type ConsumerFactory interface {
	Create() (Consumer, error)
}

// Consumer interface which every specific consumer has to implement.
type Consumer interface {
	Consume(context.Context) (<-chan Message, <-chan error, error)
	Close() error
	Info() map[string]interface{}
}

// DetermineDecoder determines the decoder based on the content type.
func DetermineDecoder(contentType string) (encoding.DecodeRawFunc, error) {
	switch contentType {
	case json.Type, json.TypeCharset:
		return json.DecodeRaw, nil
	}
	return nil, errors.Errorf("accept header %s is unsupported", contentType)
}
