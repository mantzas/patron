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
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"
)

const (
	shutdownTimeout = 5 * time.Second
)

// Component interface for implementing components.
type Component interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// Service definition.
type Service struct {
	name   string
	cps    []Component
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new service
func New(name string, cps []Component, oo ...Option) (*Service, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if len(cps) == 0 {
		return nil, errors.New("components not provided")
	}

	log.AppendField("srv", name)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log.AppendField("host", hostname)

	ctx, cancel := context.WithCancel(context.Background())
	s := Service{name: name, cps: cps, ctx: ctx, cancel: cancel}

	for _, o := range oo {
		err := o(&s)
		if err != nil {
			return nil, err
		}
	}

	s.setupTermSignal()
	return &s, nil
}

func (s *Service) setupTermSignal() {
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Info("term signal received, cancelling")
		s.cancel()
	}()
}

// Run starts up all service components and monitors for errors.
func (s *Service) Run() error {

	errCh := make(chan error)

	for _, cp := range s.cps {
		go func(c Component, ctx context.Context) {
			errCh <- c.Run(ctx)
		}(cp, s.ctx)
	}

	select {
	case err := <-errCh:
		log.Error("component returned a error")
		err1 := s.Shutdown()
		if err1 != nil {
			return errors.Wrapf(err, "failed to shutdown %v", err1)
		}
		return err
	case <-s.ctx.Done():
		log.Info("stop signal received")
		return s.Shutdown()
	}
}

// Shutdown performs a shutdown on all components with the setup timeout.
func (s *Service) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	defer func() {
		err := trace.Close()
		if err != nil {
			log.Errorf("failed to close trace %v", err)
		}
	}()
	log.Info("shutting down components")

	wg := sync.WaitGroup{}
	agr := agr_errors.New()
	for _, cp := range s.cps {

		wg.Add(1)
		go func(c Component, ctx context.Context, w *sync.WaitGroup, agr *agr_errors.Aggregate) {
			defer w.Done()
			agr.Append(c.Shutdown(ctx))
		}(cp, ctx, &wg, agr)
	}

	wg.Wait()
	if agr.Count() > 0 {
		return agr
	}
	return nil
}
