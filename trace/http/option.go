package http

import (
	"time"

	"github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/reliability/circuitbreaker"
)

// OptionFunc definition for configuring the client in a functional way.
type OptionFunc func(*TracedClient) error

// Timeout option for adjusting the timeout of the connection.
func Timeout(timeout time.Duration) OptionFunc {
	return func(tc *TracedClient) error {
		if timeout <= 0 {
			return errors.New("timeout must be positive")
		}
		tc.cl.Timeout = timeout
		return nil
	}
}

// CircuitBreaker option for setting up a circuit breaker.
func CircuitBreaker(name string, set circuitbreaker.Setting) OptionFunc {
	return func(tc *TracedClient) error {
		cb, err := circuitbreaker.New(name, set)
		if err != nil {
			return errors.Wrap(err, "failed to set circuit breaker")
		}
		tc.cb = cb
		return nil
	}
}
