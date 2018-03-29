package http

import (
	"errors"

	"github.com/mantzas/patron/log"
)

// Option defines a option for the HTTP service
type Option func(Service) error

// Ports option for setting the ports of the service and pprof
func Ports(port, pprofPort int) Option {
	return func(s Service) error {
		if port <= 0 || port > 65535 {
			return errors.New("invalid port")
		}

		if pprofPort <= 0 || pprofPort > 65535 {
			return errors.New("invalid pprof port")
		}

		if port == pprofPort {
			return errors.New("pprof must be on a separate port")

		}

		s.port = port
		s.pprofPort = pprofPort
		log.Infof("setup port: %d and pprof port: %d", port, pprofPort)
		return nil
	}
}
