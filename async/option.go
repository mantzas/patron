package async

import (
	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
)

// FailStrategy type definition.
type FailStrategy int

const (
	// NackExitStrategy does not acknowledge the message and exits the application on error.
	NackExitStrategy FailStrategy = 0
	// NackStrategy does not acknowledge the message, leaving it for reprocessing, and continues.
	NackStrategy FailStrategy = 1
	// AckStrategy acknowledges message and continues.
	AckStrategy FailStrategy = 2
)

// OptionFunc definition for configuring the component in a functional way.
type OptionFunc func(*Component) error

// FailureStrategy option for setting the strategy of handling failures in the async component.
func FailureStrategy(fs FailStrategy) OptionFunc {
	return func(c *Component) error {
		if fs > AckStrategy || fs < NackExitStrategy {
			return errors.New("invalid strategy provided")
		}
		c.failStrategy = fs
		log.Info("failure strategy set")
		return nil
	}
}
