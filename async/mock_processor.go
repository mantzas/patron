package async

import "context"

// MockMesssageProcessor definition for test usage
type MockMesssageProcessor struct {
}

// Process a message for testing purposes
func (mmp MockMesssageProcessor) Process(ctx context.Context, msg []byte) error {

	return nil
}
