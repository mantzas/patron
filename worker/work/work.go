package work

// Work interface for implementing work items
type Work interface {
}

// Acquirer interface for implementing work acquiring
type Acquirer interface {
	Acquire() ([]Work, error)
}

// Acknowledger interface for implementing work acknowledgement
type Acknowledger interface {
	Acknowledge([]Work) error
}

// Processor interface for implementing work processing
type Processor interface {
	Process([]Work) error
}
