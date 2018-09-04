package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/patron/log"
)

const (
	httpPort         = 50000
	httpReadTimeout  = 5 * time.Second
	httpWriteTimeout = 10 * time.Second
	httpIdleTimeout  = 120 * time.Second
)

var (
	// DefaultHealthCheck returns always healthy.
	DefaultHealthCheck = func() HealthStatus { return Healthy }
)

// Component implementation of HTTP.
type Component struct {
	hc               HealthCheckFunc
	httpPort         int
	httpReadTimeout  time.Duration
	httpWriteTimeout time.Duration
	sync.Mutex
	routes   []Route
	certFile string
	keyFile  string
}

// New returns a new component.
func New(oo ...OptionFunc) (*Component, error) {
	s := Component{
		hc:               DefaultHealthCheck,
		httpPort:         httpPort,
		httpReadTimeout:  httpReadTimeout,
		httpWriteTimeout: httpWriteTimeout,
		routes:           []Route{},
	}

	for _, o := range oo {
		err := o(&s)
		if err != nil {
			return nil, err
		}
	}

	s.routes = append(s.routes, healthCheckRoute(s.hc))
	s.routes = append(s.routes, profilingRoutes()...)
	s.routes = append(s.routes, metricRoute())

	return &s, nil
}

// Run starts the HTTP server.
func (s *Component) Run(ctx context.Context) error {
	s.Lock()
	log.Debug("applying tracing to routes")
	for i := 0; i < len(s.routes); i++ {
		if s.routes[i].Trace {
			s.routes[i].Handler = DefaultMiddleware(s.routes[i].Pattern, s.routes[i].Handler)
		} else {
			s.routes[i].Handler = RecoveryMiddleware(s.routes[i].Handler)
		}
	}
	chFail := make(chan error)
	srv := s.createHTTPServer()
	go s.listenAndServe(srv, chFail)
	s.Unlock()

	select {
	case <-ctx.Done():
		log.Info("shutting down component")
		return srv.Shutdown(ctx)
	case err := <-chFail:
		return err
	}
}

func (s *Component) listenAndServe(srv *http.Server, ch chan<- error) {
	if s.certFile != "" && s.keyFile != "" {
		log.Infof("HTTPS component listening on port %d", s.httpPort)
		ch <- srv.ListenAndServeTLS(s.certFile, s.keyFile)
	}

	log.Infof("HTTP component listening on port %d", s.httpPort)
	ch <- srv.ListenAndServe()
}

func (s *Component) createHTTPServer() *http.Server {
	log.Debugf("adding %d routes", len(s.routes))
	router := httprouter.New()
	for _, route := range s.routes {
		router.HandlerFunc(route.Method, route.Pattern, route.Handler)
		log.Debugf("added route %s %s", route.Method, route.Pattern)
	}
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.httpPort),
		ReadTimeout:  s.httpReadTimeout,
		WriteTimeout: s.httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
		Handler:      router,
	}
}
