package patron

import (
	"context"
	"net/http"
	_ "net/http/pprof" // Package blank import for supporting pprof
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Service definition
type Service struct {
	pprofSrv *http.Server
}

// New constructs a new service with functional configuration.
func New(options ...Option) (*Service, error) {

	s := &Service{
		pprofSrv: &http.Server{
			Addr:         ":81",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      http.DefaultServeMux,
		},
	}

	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// ListenAndServe starts up the server, listens and serves requests
func (s *Service) ListenAndServe() error {
	return s.pprofSrv.ListenAndServe()
}

// WaitSignalAndShutdown awaits a SIGTERM to shut down
func (s *Service) WaitSignalAndShutdown(timeout time.Duration) error {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.pprofSrv.Shutdown(ctx)
}
