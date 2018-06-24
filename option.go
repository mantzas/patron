package patron

import (
	"errors"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/sync/http"
)

// Option defines a option for the HTTP service.
type Option func(*Service) error

// Routes option for adding routes to the default HTTP component.
func Routes(rr []http.Route) Option {
	return func(s *Service) error {
		if rr == nil || len(rr) == 0 {
			return errors.New("routes are required")
		}
		s.routes = rr
		log.Info("routes options are set")
		return nil
	}
}

// HealthCheck option for adding a custom health check to the default HTTP component.
func HealthCheck(hcf http.HealthCheckFunc) Option {
	return func(s *Service) error {
		if hcf == nil {
			return errors.New("health check func is required")
		}
		s.hcf = hcf
		log.Info("health check func is set")
		return nil
	}
}

// Components option for adding additional components to the service.
func Components(cc ...Component) Option {
	return func(s *Service) error {
		if cc == nil || len(cc) == 0 || cc[0] == nil {
			return errors.New("components are required")
		}
		s.cps = append(s.cps, cc...)
		log.Info("component options are set")
		return nil
	}
}
