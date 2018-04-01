package worker

// MessageProcessor interface for implementing processing of messages
type MessageProcessor interface {
	Process([]byte) error
}

// Processor interface to be implemented by components that get data from message brokers etc
type Processor interface {
	Process() error
}
