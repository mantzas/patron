package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

type handlerGen func([]Route) http.Handler

const (
	port = 50000
)

var (
	defaultHealthCheck = func() HealthStatus { return Healthy }
)

// Component implementation of HTTP.
type Component struct {
	hg     handlerGen
	hc     HealthCheckFunc
	port   int
	routes []Route
	m      sync.Mutex
	srv    *http.Server
}

// New returns a new component.
func New(hg handlerGen, oo ...Option) (*Component, error) {
	if hg == nil {
		return nil, errors.New("http handler generator is required")
	}

	s := Component{hg, defaultHealthCheck, port, []Route{}, sync.Mutex{}, nil}

	for _, o := range oo {
		err := o(&s)
		if err != nil {
			return nil, err
		}
	}

	s.routes = append(s.routes, healthCheckRoute(s.hc))
	s.routes = append(s.routes, profilingRoutes()...)

	return &s, nil
}

// Run starts the HTTP server.
func (s *Component) Run(ctx context.Context) error {
	s.m.Lock()
	log.Infof("applying tracing to routes")
	for i := 0; i < len(s.routes); i++ {
		if s.routes[i].Trace {
			s.routes[i].Handler = DefaultMiddleware(s.routes[i].Pattern, s.routes[i].Handler)
		} else {
			s.routes[i].Handler = RecoveryMiddleware(s.routes[i].Handler)
		}
	}
	s.srv = createHTTPServer(s.port, s.hg(s.routes))
	s.m.Unlock()
	log.Infof("component listening on port %d", s.port)
	return s.srv.ListenAndServe()
}

// Shutdown the component.
func (s *Component) Shutdown(ctx context.Context) error {
	s.m.Lock()
	defer s.m.Unlock()
	log.Info("shutting down component")
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
