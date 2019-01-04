package confluent

import (
	"errors"
)

// OptionFunc definition for configuring the consumer in a functional way.
type OptionFunc func(*consumer) error

// Config option for configuring consumer.
func Config(cfg map[string]interface{}) OptionFunc {
	return func(c *consumer) error {
		if cfg == nil {
			return errors.New("config is nil")
		}

		if len(cfg) == 0 {
			return errors.New("config is empty")
		}

		for k, v := range cfg {
			c.cfg.SetKey(k, v)
		}
		return nil
	}
}
