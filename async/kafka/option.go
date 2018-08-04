package kafka

import (
	"time"

	"github.com/pkg/errors"
)

// OptionFunc definition for configuring the consumer in a functional way.
type OptionFunc func(*Consumer) error

//buffer int, start Offset,

// Buffer option for adjusting the incoming messages buffer.
func Buffer(buf int) OptionFunc {
	return func(c *Consumer) error {
		if buf < 0 {
			return errors.New("buffer must greater or equal than 0")
		}
		c.buffer = buf
		return nil
	}
}

// Start option for adjusting the start point in the topic.
func Start(start Offset) OptionFunc {
	return func(c *Consumer) error {
		c.start = start
		return nil
	}
}

// Timeout option for adjusting the timeout of the connection.
func Timeout(timeout time.Duration) OptionFunc {
	return func(c *Consumer) error {
		c.cfg.Net.DialTimeout = timeout
		return nil
	}
}
