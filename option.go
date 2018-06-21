package patron

import (
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"
	jaeger "github.com/uber/jaeger-client-go"
)

// Option defines a option for the HTTP service.
type Option func(*Service) error

// Tracing option for setting tracing.
func Tracing(agentAddress, samplerType string, samplerParam float64) Option {
	return func(s *Service) error {
		if agentAddress == "" {
			return errors.New("agent address is required")
		}
		if samplerType != jaeger.SamplerTypeConst &&
			samplerType != jaeger.SamplerTypeRemote &&
			samplerType != jaeger.SamplerTypeProbabilistic &&
			samplerType != jaeger.SamplerTypeRateLimiting &&
			samplerType != jaeger.SamplerTypeLowerBound {
			return errors.New("invalid sampler type provided")
		}
		err := trace.Setup(s.name, agentAddress, samplerType, samplerParam)
		if err != nil {
			return err
		}
		log.Info("tracing set")
		return nil
	}
}
