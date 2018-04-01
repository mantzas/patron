package worker

import (
	"os"

	"github.com/mantzas/patron"

	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

// Service definition
type Service struct {
	patron.Service
	p Processor
}

// New creates a new service
func New(name string, p Processor) (*Service, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if p == nil {
		return nil, errors.New("processor is required")
	}

	log.AppendField("wrk", name)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log.AppendField("host", hostname)

	return &Service{*patron.New(), p}, nil
}

// Run kicks off the processing
func (s Service) Run() error {
	return errors.Wrap(s.p.Process(s.Ctx), "processor failed")
}
