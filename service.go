package patron

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/mantzas/patron/sync/http"

	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
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
	routes []http.Route
	hcf    http.HealthCheckFunc
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new service
func New(name string, oo ...Option) (*Service, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	err := setupDefaultLogging(name)
	if err != nil {
		return nil, err
	}

	err = setupDefaultTracing(name)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := Service{name: name, cps: []Component{}, hcf: http.DefaultHealthCheck, ctx: ctx, cancel: cancel}

	for _, o := range oo {
		err := o(&s)
		if err != nil {
			return nil, err
		}
	}

	httpCp, err := s.createHTTPComponent()
	if err != nil {
		return nil, err
	}

	s.cps = append(s.cps, httpCp)
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

func setupDefaultLogging(srvName string) error {
	lvl, ok := os.LookupEnv("PATRON_LOG_LEVEL")
	if !ok {
		lvl = string(log.InfoLevel)
	}

	err := log.Setup(zerolog.DefaultFactory(log.Level(lvl)))
	if err != nil {
		return errors.Wrap(err, "failed to setup logging")
	}

	log.AppendField("srv", srvName)
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "failed to get hostname")
	}
	log.AppendField("host", hostname)
	log.Info("set up default log level to `INFO`")
	return nil
}

func setupDefaultTracing(srvName string) error {
	agent, ok := os.LookupEnv("PATRON_JAEGER_AGENT")
	if !ok {
		agent = "0.0.0.0:6831"
	}
	tp, ok := os.LookupEnv("PATRON_JAEGER_SAMPLER_TYPE")
	if !ok {
		tp = jaeger.SamplerTypeProbabilistic
	}
	prm, ok := os.LookupEnv("PATRON_JAEGER_SAMPLER_PARAM")
	if !ok {
		prm = "0.1"
	}

	param, err := strconv.ParseFloat(prm, 64)
	if err != nil {
		return errors.Wrap(err, "failed to convet sampler param to float64")
	}

	log.Infof("setting up default tracing to %s, %s with param %s", agent, tp, prm)
	return trace.Setup(srvName, agent, tp, param)
}

func (s *Service) createHTTPComponent() (Component, error) {

	port, ok := os.LookupEnv("PATRON_HTTP_DEFAULT_PORT")
	if !ok {
		port = "50000"
	}

	log.Infof("creating default HTTP component at port %s", port)

	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("failed to parse port %s", port)
	}

	options := []http.Option{
		http.Port(p),
	}

	if s.hcf != nil {
		options = append(options, http.HealthCheck(s.hcf))
	}

	if s.routes != nil {
		options = append(options, http.Routes(s.routes))
	}

	cp, err := http.New(options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create default HTTP component")
	}

	return cp, nil
}
