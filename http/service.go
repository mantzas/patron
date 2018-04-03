package http

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/http/pprof"
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

const (
	port            = 50000
	pprofPort       = 50001
	shutdownTimeout = 5 * time.Second
)

// Service definition for handling HTTP request
type Service struct {
	patron.Service
	port       int
	pprofPort  int
	HandlerGen HandlerGen
	srv        *http.Server
	pprof      *pprof.Server
}

// New returns a new service with options applied
func New(name string, routes []Route, options ...Option) (*Service, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if len(routes) == 0 {
		return nil, errors.New("routes should be provided")
	}

	s := Service{
		Service:    *patron.New(),
		port:       port,
		pprofPort:  pprofPort,
		HandlerGen: CreateHandler,
	}

	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}
	log.AppendField("srv", name)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log.AppendField("host", hostname)
	s.srv = CreateHTTPServer(s.port, s.HandlerGen(routes))
	s.pprof = pprof.New(s.pprofPort)

	return &s, nil
}

// Run starts up the server, listens and serves requests
func (s *Service) Run() error {

	errCh := make(chan error)

	go func() {
		log.Infof("listen and server pprof on port %d", s.pprofPort)
		errCh <- s.pprof.ListenAndServe()
	}()

	go func() {
		log.Infof("listen and server service on port %d", s.port)
		errCh <- s.srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		log.Info("service/pprof returned a error")
		err1 := s.shutdown()
		if err1 != nil {
			return errors.Wrapf(err, "failed to shutdown %v", err1)
		}
		return err
	case <-s.Ctx.Done():
		log.Info("stop signal received")
		return s.shutdown()
	}
}

func (s *Service) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Info("shutting down pprof")
	err := s.pprof.Shutdown(ctx)
	if err != nil {
		log.Error("failed to shutdown pprof server")
	}

	log.Info("shutting down service")
	return s.srv.Shutdown(ctx)
}
