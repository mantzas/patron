package patron

import (
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"
)

// Option defines a option for the HTTP service.
type Option func(*Service) error

// Tracing option for setting tracing.
func Tracing(agentAddress string) Option {
	return func(s *Service) error {
		if agentAddress == "" {
			return errors.New("agent address is required")
		}
		err := trace.Setup(s.name, agentAddress)
		if err != nil {
			return err
		}
		log.Info("tracing set")
		return nil
	}
}
