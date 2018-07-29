package patron

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mantzas/patron/config"
	"github.com/mantzas/patron/config/env"
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

// Component interface for implementing service components.
type Component interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// Service is responsible for managing and setting up everything.
// The service will start by default a HTTP component in order to host management endpoint.
type Service struct {
	name   string
	cps    []Component
	routes []http.Route
	hcf    http.HealthCheckFunc
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new named service and allows for customization through functional options.
func New(name, version string, oo ...OptionFunc) (*Service, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if version == "" {
		version = "dev"
	}

	err := setupDefaultConfig()
	if err != nil {
		return nil, err
	}

	err = setupDefaultLogging(name, version)
	if err != nil {
		return nil, err
	}

	err = setupDefaultTracing(name, version)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := Service{name: name, cps: []Component{}, hcf: http.DefaultHealthCheck, ctx: ctx, cancel: cancel}

	for _, o := range oo {
		err = o(&s)
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
// If a component returns a error the service is responsible for shutting down
// all components and terminate itself.
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

// Shutdown all components gracefully with a predefined timeout.
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

func setupDefaultConfig() error {
	f, err := os.Open(".env")
	if err != nil {
		f = nil
	}

	cfg, err := env.New(f)
	if err != nil {
		return err
	}

	return config.Setup(cfg)
}

func setupDefaultLogging(name, version string) error {
	lvl, err := config.GetString("PATRON_LOG_LEVEL")
	if err != nil {
		lvl = string(log.InfoLevel)
	}

	err = log.Setup(zerolog.DefaultFactory(log.Level(lvl)))
	if err != nil {
		return errors.Wrap(err, "failed to setup logging")
	}

	log.AppendField("srv", name)
	log.AppendField("version", version)
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "failed to get hostname")
	}
	log.AppendField("host", hostname)
	log.Info("set up default log level to `INFO`")
	return nil
}

func setupDefaultTracing(name, version string) error {
	agent, err := config.GetString("PATRON_JAEGER_AGENT")
	if err != nil {
		agent = "0.0.0.0:6831"
	}
	tp, err := config.GetString("PATRON_JAEGER_SAMPLER_TYPE")
	if err != nil {
		tp = jaeger.SamplerTypeProbabilistic
	}
	prm, err := config.GetFloat64("PATRON_JAEGER_SAMPLER_PARAM")
	if err != nil {
		prm = 0.1
	}
	log.Infof("setting up default tracing to %s, %s with param %f", agent, tp, prm)
	return trace.Setup(name, version, agent, tp, prm)
}

func (s *Service) createHTTPComponent() (Component, error) {

	port, err := config.GetInt64("PATRON_HTTP_DEFAULT_PORT")
	if err != nil {
		port = 50000
	}

	log.Infof("creating default HTTP component at port %d", port)

	options := []http.OptionFunc{
		http.Port(int(port)),
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
