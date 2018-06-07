package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mantzas/patron/log"
	"github.com/opentracing/opentracing-go"
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
	tr     opentracing.Tracer
	hg     handlerGen
	hc     HealthCheckFunc
	port   int
	routes []Route
	srv    *http.Server
}

// New returns a new component.
func New(tr opentracing.Tracer, hg handlerGen, oo ...Option) (*Component, error) {
	if tr == nil {
		return nil, errors.New("tracer is required")
	}

	if hg == nil {
		return nil, errors.New("http handler generator is required")
	}

	s := Component{tr, hg, defaultHealthCheck, port, []Route{}, nil}

	for _, o := range oo {
		err := o(&s)
		if err != nil {
			return nil, err
		}
	}

	s.routes = append(s.routes, healthCheckRoute(s.hc))
	s.routes = append(s.routes, profilingRoutes()...)

	s.srv = createHTTPServer(s.port, s.hg(s.routes))
	return &s, nil
}

// Run starts the HTTP server.
func (s *Component) Run(ctx context.Context) error {
	log.Infof("component listening on port %d", s.port)
	return s.srv.ListenAndServe()
}

// Shutdown the component.
func (s *Component) Shutdown(ctx context.Context) error {
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
