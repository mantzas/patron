package patron

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

const shutdownTimeout = 5 * time.Second

type Service interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// Server definition of a server hosting service
type Server struct {
	name     string
	services []Service
	Ctx      context.Context
	Cancel   context.CancelFunc
}

// New creates a new server
func New(name string, services ...Service) (*Server, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if len(services) == 0 {
		return nil, errors.New("services not provided")
	}

	log.AppendField("srv", name)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log.AppendField("host", hostname)

	ctx, cancel := context.WithCancel(context.Background())
	s := Server{name, services, ctx, cancel}

	s.setupTermSignal()
	return &s, nil
}

func (s *Server) setupTermSignal() {
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Info("term signal received, cancelling")
		s.Cancel()
	}()
}

// Run starts up the server, listens and serves requests
func (s *Server) Run() error {

	errCh := make(chan error)

	for _, service := range s.services {
		go func(s Service, ctx context.Context) {
			errCh <- s.Run(ctx)
		}(service, s.Ctx)
	}

	select {
	case err := <-errCh:
		log.Error("service returned a error")
		err1 := s.Shutdown()
		if err1 != nil {
			return errors.Wrapf(err, "failed to shutdown %v", err1)
		}
		return err
	case <-s.Ctx.Done():
		log.Info("stop signal received")
		return s.Shutdown()
	}
}

// Shutdown performs a shutdown on all services with the setup timeout
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	log.Info("shutting down services")

	wg := sync.WaitGroup{}
	wg.Add(len(s.services))
	agr := agr_errors.New()

	for _, srv := range s.services {

		go func(srv Service, ctx context.Context, w *sync.WaitGroup, agr *agr_errors.Aggregate) {
			defer w.Done()
			agr.Append(srv.Shutdown(ctx))
		}(srv, ctx, &wg, agr)
	}

	wg.Wait()
	if agr.Count() > 0 {
		return agr
	}
	return nil
}
