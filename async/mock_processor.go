package async

import "context"

// MockProcessor definition for test usage.
type MockProcessor struct {
}

// Process a message for testing purposes.
func (mp MockProcessor) Process(ctx context.Context, msg *Message) error {
	return nil
}
