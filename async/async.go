package async

import (
	"context"

	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/pkg/errors"
)

const (
	// ContentTypeHeader constant
	ContentTypeHeader string = "Content-Type"
)

// Processor interface for implementing processing of messages
type Processor interface {
	Process(context.Context, *Message) error
}

// DetermineDecoder determines the decoder based on the content type.
func DetermineDecoder(contentType string) (encoding.DecodeRaw, error) {
	switch contentType {
	case json.ContentType, json.ContentTypeCharset:
		return json.DecodeRaw, nil
	}
	return nil, errors.Errorf("accept header %s is unsupported", contentType)
}
