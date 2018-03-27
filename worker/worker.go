package worker

import (
	"errors"
	"os"

	"github.com/mantzas/patron/log"
)

// Service definition
type Service struct {
}

// New creates a new service
func New(name string) (*Service, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	log.AppendField("srv", name)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log.AppendField("host", hostname)
	log.Info("creating a new service")

	return &Service{}, nil
}

// Run kicks off the processing
func (s Service) Run() error {

	return nil
}
