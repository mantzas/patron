package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mantzas/patron/http/pprof"
)

// Service definition for handling HTTP request
type Service struct {
	srv   *http.Server
	pprof *pprof.Server
}

// New returns a new service with options applied
func New(options ...Option) (*Service, error) {

	// TODO: replace with actual mux
	s := Service{
		srv:   CreateHTTPServer(80, http.DefaultServeMux),
		pprof: pprof.New(81),
	}

	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	return &s, nil
}

// ListenAndServe starts up the server, listens and serves requests
func (s *Service) ListenAndServe() error {

	errCh := make(chan error)

	go func() {
		errCh <- s.pprof.ListenAndServe()
	}()

	go func() {
		errCh <- s.srv.ListenAndServe()
	}()

	return <-errCh
}

// WaitSignalAndShutdown waiting for shutdown signal and shuts down service
func (s *Service) WaitSignalAndShutdown(timeout time.Duration) error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := s.pprof.Shutdown(ctx)
	if err != nil {
		//TODO: logging
	}

	return s.srv.Shutdown(ctx)
}
