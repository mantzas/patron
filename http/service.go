package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mantzas/patron/http/pprof"
	"github.com/pkg/errors"
)

const (
	port            = 50000
	pprofPort       = 50001
	shutdownTimeout = 5 * time.Second
)

// Service definition for handling HTTP request
type Service struct {
	srv   *http.Server
	pprof *pprof.Server
}

// New returns a new service with options applied
func New(options ...Option) (*Service, error) {

	// TODO: replace with actual mux
	mux := http.ServeMux{}
	s := Service{
		srv:   CreateHTTPServer(port, &mux),
		pprof: pprof.New(pprofPort),
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:

		err1 := s.shutdown()
		if err1 != nil {
			return errors.Wrapf(err, "failed to shutdown %v", err1)
		}
		return err

	case <-stop:
		return s.shutdown()
	}
}

func (s *Service) shutdown() error {

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := s.pprof.Shutdown(ctx)
	if err != nil {
		//TODO: logging
	}

	return s.srv.Shutdown(ctx)
}
