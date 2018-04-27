package http

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/mantzas/patron/log"
)

const (
	port = 50000
)

var (
	defaultHealthCheck = func() HealthStatus { return Healthy }
)

// Service implementation of HTTP
type Service struct {
	hg     handlerGen
	hc     HealthCheckFunc
	port   int
	routes []Route
	srv    *http.Server
	m      sync.Mutex
}

// New returns a new service
func New(hg handlerGen, options ...Option) (*Service, error) {
	if hg == nil {
		return nil, errors.New("http handler generator is required")
	}

	s := Service{hg, defaultHealthCheck, port, []Route{}, nil, sync.Mutex{}}

	for _, opt := range options {
		err := opt(&s)
		if err != nil {
			return nil, err
		}
	}

	s.routes = append(s.routes, healthCheckRoute(s.hc))
	s.srv = createHTTPServer(s.port, s.hg(s.routes))
	return &s, nil
}

// Run starts the processing
func (s *Service) Run(ctx context.Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	log.Infof("service listening on port %d", s.port)
	return s.srv.ListenAndServe()
}

// Shutdown the service
func (s *Service) Shutdown(ctx context.Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	log.Info("shutting down service")
	return s.srv.Shutdown(ctx)
}
