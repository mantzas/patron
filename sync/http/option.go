package http

import (
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

// Option defines a option for the HTTP service
type Option func(*Service) error

// SetPorts option for setting the ports of the http service
func SetPorts(port int) Option {
	return func(s *Service) error {
		if port <= 0 || port > 65535 {
			return errors.New("invalid port")
		}
		s.port = port
		log.Infof("port set to %d", port)
		return nil
	}
}

// SetRoutes option for setting the routes of the http service
func SetRoutes(rr []Route) Option {
	return func(s *Service) error {
		if len(rr) == 0 {
			return errors.New("routes are empty")
		}
		s.routes = rr
		log.Info("routes set")
		return nil
	}
}
