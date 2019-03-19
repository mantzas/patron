package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/errors"
)

// OptionFunc definition for configuring the consumer in a functional way.
type OptionFunc func(*consumer) error

// Version option for setting the Kafka version.
func Version(version string) OptionFunc {
	return func(c *consumer) error {
		if version == "" {
			return errors.New("versions has to be provided")
		}

		v, err := sarama.ParseKafkaVersion(version)
		if err != nil {
			return errors.Wrap(err, "invalid kafka version provided")
		}

		c.cfg.Version = v
		return nil
	}
}

// Buffer option for adjusting the incoming messages buffer.
func Buffer(buf int) OptionFunc {
	return func(c *consumer) error {
		if buf < 0 {
			return errors.New("buffer must greater or equal than 0")
		}
		c.buffer = buf
		return nil
	}
}

// Start option for adjusting the start point in the topic.
func Start(start Offset) OptionFunc {
	return func(c *consumer) error {
		c.start = start
		return nil
	}
}

// Timeout option for adjusting the timeout of the connection.
func Timeout(timeout time.Duration) OptionFunc {
	return func(c *consumer) error {
		c.cfg.Net.DialTimeout = timeout
		return nil
	}
}
