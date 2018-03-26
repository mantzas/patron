package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	srv   *http.Server
	pprof *pprof.Server
}

// New returns a new service with options applied
func New(name string, options ...Option) (*Service, error) {

	log.AppendField("srv", name)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log.AppendField("host", hostname)
	log.Info("creating a new service")
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
		log.Infof("listen and server pprof:%s", s.pprof.GetAddr())
		errCh <- s.pprof.ListenAndServe()
	}()

	go func() {
		log.Infof("listen and server service:%s", s.srv.Addr)
		errCh <- s.srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Info("service/pprof returned a error")
		err1 := s.shutdown()
		if err1 != nil {
			return errors.Wrapf(err, "failed to shutdown %v", err1)
		}
		return err
	case <-stop:
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
