package kafka

import (
	"errors"
	"fmt"
)

// OptionFunc definition for configuring the consumer in a functional way.
type OptionFunc func(*consumer) error

// Config option for configuring consumer.
// See https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md.
func Config(cfg map[string]interface{}) OptionFunc {
	return func(c *consumer) error {
		if cfg == nil {
			return errors.New("config is nil")
		}

		if len(cfg) == 0 {
			return errors.New("config is empty")
		}

		var err error

		for k, v := range cfg {
			err = c.cfg.SetKey(k, v)
			if err != nil {
				return fmt.Errorf("failed to set key %s: %v", k, err)
			}
			if k == "enable.auto.commit" {
				auto, ok := v.(bool)
				if !ok {
					return errors.New("enable.auto.commit should be boolean")
				}
				c.ack = !auto
			}
		}
		return nil
	}
}
