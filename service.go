package patron

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

const (
	port            = 50000
	shutdownTimeout = 5 * time.Second
)

// Service base component
type Service struct {
	name           string
	port           int
	routes         []Route
	HTTPHandlerGen httpHandlerGen
	proc           Processor
	srv            *http.Server
	Ctx            context.Context
	Cancel         context.CancelFunc
}

// New creates a new base service
func New(name string, httpHandlerGen httpHandlerGen, options ...Option) (*Service, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if httpHandlerGen == nil {
		return nil, errors.New("http handler generator is required")
	}

	log.AppendField("srv", name)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	log.AppendField("host", hostname)

	ctx, cancel := context.WithCancel(context.Background())
	s := Service{
		name:           name,
		port:           port,
		routes:         []Route{},
		Ctx:            ctx,
		Cancel:         cancel,
		proc:           nil,
		HTTPHandlerGen: httpHandlerGen,
	}

	for _, opt := range options {
		err := opt(&s)
		if err != nil {
			return nil, err
		}
	}

	s.setupTermSignal()
	s.srv = createHTTPServer(s.port, s.HTTPHandlerGen(s.routes))
	return &s, nil
}

func (s *Service) setupTermSignal() {
	go func() {

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Info("term signal received, cancelling")
		s.Cancel()
	}()
}

// Run starts up the server, listens and serves requests
func (s *Service) Run() error {

	errCh := make(chan error)

	go func() {
		log.Infof("service listening on port %d", s.port)
		errCh <- s.srv.ListenAndServe()
	}()

	if s.proc != nil {
		go func() {
			log.Info("starting processing")
			errCh <- s.proc.Process(s.Ctx)
		}()
	}

	select {
	case err := <-errCh:
		log.Info("service returned a error")
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

	log.Info("shutting down service")
	return s.srv.Shutdown(ctx)
}
