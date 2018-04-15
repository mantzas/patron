package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/mantzas/patron/log"
)

const (
	port = 50000
)

// Service implementation of HTTP
type Service struct {
	hg     handlerGen
	port   int
	routes []Route
	srv    *http.Server
}

// New returns a new service
func New(hg handlerGen, options ...Option) (*Service, error) {
	if hg == nil {
		return nil, errors.New("http handler generator is required")
	}

	s := Service{hg, port, []Route{}, nil}

	for _, opt := range options {
		err := opt(&s)
		if err != nil {
			return nil, err
		}
	}

	s.srv = createHTTPServer(s.port, s.hg(s.routes))
	return &s, nil
}

// Run starts the processing
func (s *Service) Run(ctx context.Context) error {

	log.Infof("service listening on port %d", s.port)
	return s.srv.ListenAndServe()
}

// Shutdown the service
func (s *Service) Shutdown(ctx context.Context) error {

	log.Info("shutting down service")
	return s.srv.Shutdown(ctx)
}
