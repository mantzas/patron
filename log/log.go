package log

import (
	"errors"
	"os"
)

var factory Factory

func init() {
	factory = NewStdFactory(os.Stdout)
}

// Setup set's up a new factory to the global state
func Setup(f Factory) error {
	if f == nil {
		return errors.New("factory is nil")
	}

	factory = f

	return nil
}

// Create returns a new logger. Fields are optional and allow nil
func Create(f map[string]interface{}) Logger {
	return factory.Create(f)
}

// CreateSub returns a new sub logger
func CreateSub(l Logger, fields map[string]interface{}) Logger {
	return factory.CreateSub(l, fields)
}
