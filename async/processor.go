package async

import "context"

// Processor interface for implementing processing of messages
type Processor interface {
	Process(context.Context, []byte) error
}
