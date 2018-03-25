package patron

import (
	"errors"
	"fmt"

	"github.com/mantzas/patron/config"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/metric"
)

// Option defines a option for the service
type Option func(*Service) error

// Config option for setting the configuration of the service
func Config(c config.Config) Option {
	return func(p *Service) error {
		return config.Setup(c)
	}
}

// Log option for setting the logging of the service
func Log(f log.Factory) Option {
	return func(p *Service) error {
		return log.Setup(f)
	}
}

// Metric option for setting the metrics of the service
func Metric(m metric.Metric) Option {
	return func(p *Service) error {
		return metric.Setup(m)
	}
}

// PProf option for setting the port of pprof of the service
func PProf(port int) Option {
	return func(p *Service) error {
		if port <= 0 || port > 65535 {
			return errors.New("invalid pprof port")
		}

		p.pprofSrv.Addr = fmt.Sprintf(":%d", port)
		return nil
	}
}
