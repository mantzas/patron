package patron

import (
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
)

// Option defines a option for the HTTP service.
type Option func(*Service) error

// Tracing option for setting tracing.
func Tracing(sampler jaeger.Sampler, reporter jaeger.Reporter, options ...jaeger.TracerOption) Option {
	return func(s *Service) error {
		if sampler == nil {
			return errors.New("sampler is required")
		}
		if reporter == nil {
			return errors.New("reporter is required")
		}

		tr, trCloser := jaeger.NewTracer(s.name, sampler, reporter, options...)
		s.tr = tr
		s.trCloser = trCloser
		log.Info("tracing set")
		return nil
	}
}
