package worker

import (
	"os"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/worker/work"
	"github.com/pkg/errors"
)

// Service definition
type Service struct {
	acq work.Acquirer
	prc work.Processor
}

// New creates a new service
func New(name string, acq work.Acquirer, prc work.Processor) (*Service, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if acq == nil {
		return nil, errors.New("acquirer is required")
	}

	if prc == nil {
		return nil, errors.New("processor is required")
	}

	log.AppendField("wrk", name)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log.AppendField("host", hostname)

	return &Service{acq, prc}, nil
}

// Run kicks off the processing
func (s Service) Run() error {

	for {
		w, err := s.acq.Acquire()
		if err != nil {
			return errors.Wrap(err, "failed to acquire work")
		}

		err = s.prc.Process(w)
		if err != nil {
			return errors.Wrap(err, "failed to process work")
		}
	}
}
