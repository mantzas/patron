package kafka

import (
	"fmt"

	"github.com/mantzas/patron/errors"
)

// OptionFunc definition for configuring the async producer in a functional way.
type OptionFunc func(*KafkaProducer) error

// Config option for configuring consumer.
func Config(cfg map[string]interface{}) OptionFunc {
	return func(kp *KafkaProducer) error {
		if cfg == nil {
			return errors.New("config is nil")
		}

		if len(cfg) == 0 {
			return errors.New("config is empty")
		}

		var err error

		for k, v := range cfg {
			err = kp.cfg.SetKey(k, v)
			if err != nil {
				return fmt.Errorf("failed to set key %s: %v", k, err)
			}
		}
		return nil
	}
}
