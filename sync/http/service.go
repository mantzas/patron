package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
	"go.opencensus.io/stats/view"
)

type handlerGen func([]Route) http.Handler

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
func New(hg handlerGen, oo ...Option) (*Service, error) {
	if hg == nil {
		return nil, errors.New("http handler generator is required")
	}

	s := Service{hg, defaultHealthCheck, port, []Route{}, nil, sync.Mutex{}}

	for _, o := range oo {
		err := o(&s)
		if err != nil {
			return nil, err
		}
	}

	s.routes = append(s.routes, healthCheckRoute(s.hc))
	s.routes = append(s.routes, profilingRoutes()...)

	s.srv = createHTTPServer(s.port, s.hg(s.routes))
	if err := view.Register(defaultServerViews...); err != nil {
		return nil, errors.Wrap(err, "failed to register default server views")
	}
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

func createHTTPServer(port int, sm http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      sm,
	}
}
