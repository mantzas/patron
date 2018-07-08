package http

import (
	"github.com/pkg/errors"
)

// OptionFunc defines a option func for the HTTP component.
type OptionFunc func(*Component) error

// Port option for setting the ports of the HTTP component.
func Port(port int) OptionFunc {
	return func(s *Component) error {
		if port <= 0 || port > 65535 {
			return errors.New("invalid port")
		}
		s.port = port
		return nil
	}
}

// Routes option for setting the routes of the HTTP component.
func Routes(rr []Route) OptionFunc {
	return func(s *Component) error {
		if len(rr) == 0 {
			return errors.New("routes are empty")
		}
		s.routes = append(s.routes, rr...)
		return nil
	}
}

// HealthCheck option for setting the health check function of the HTTP component.
func HealthCheck(hcf HealthCheckFunc) OptionFunc {
	return func(s *Component) error {
		if hcf == nil {
			return errors.New("health check function is not defined")
		}
		s.hc = hcf
		return nil
	}
}
