package worker

import "context"

// MessageProcessor interface for implementing processing of messages
type MessageProcessor interface {
	Process(context.Context, []byte) error
}

// Processor interface to be implemented by components that get data from message brokers etc
type Processor interface {
	Process(context.Context) error
}
