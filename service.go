package patron

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	agr_errors "github.com/mantzas/patron/errors"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/sync/http"
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
	log    log.Logger
}

// New creates a new named service and allows for customization through functional options.
func New(cfg Config, oo ...OptionFunc) (*Service, error) {

	if cfg.Name == "" {
		return nil, errors.New("name is required")
	}

	if cfg.Version == "" {
		cfg.Version = "dev"
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := Service{name: cfg.Name, cps: []Component{}, hcf: http.DefaultHealthCheck, ctx: ctx, cancel: cancel, log: log.Create()}

	err := s.setupDefaultTracing(cfg)
	if err != nil {
		return nil, err
	}

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
		s.log.Info("term signal received, cancelling")
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
		s.log.Error("component returned a error")
		err1 := s.Shutdown()
		if err1 != nil {
			return errors.Wrapf(err, "failed to shutdown %v", err1)
		}
		return err
	case <-s.ctx.Done():
		s.log.Info("stop signal received")
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
			s.log.Errorf("failed to close trace %v", err)
		}
	}()
	s.log.Info("shutting down components")

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

func (s *Service) setupDefaultTracing(cfg Config) error {
	agent, ok := os.LookupEnv("PATRON_JAEGER_AGENT")
	if !ok {
		agent = "0.0.0.0:6831"
	}
	tp, ok := os.LookupEnv("PATRON_JAEGER_SAMPLER_TYPE")
	if !ok {
		tp = jaeger.SamplerTypeProbabilistic
	}
	var prmVal float64
	var err error

	prm, ok := os.LookupEnv("PATRON_JAEGER_SAMPLER_PARAM")
	if !ok {
		prmVal = 0.1
	} else {
		prmVal, err = strconv.ParseFloat(prm, 64)
		if err != nil {
			return errors.Wrap(err, "env var for jaeger sampler param is not valid")
		}
	}
	s.log.Infof("setting up default tracing to %s, %s with param %f", agent, tp, prm)
	return trace.Setup(cfg.Name, cfg.Version, agent, tp, prmVal)
}

func (s *Service) createHTTPComponent() (Component, error) {
	var err error
	var portVal int64
	port, ok := os.LookupEnv("PATRON_HTTP_DEFAULT_PORT")
	if !ok {
		portVal = 50000
	} else {
		portVal, err = strconv.ParseInt(port, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "env var for HTTP default port is not valid")
		}
	}

	s.log.Infof("creating default HTTP component at port %d", port)

	options := []http.OptionFunc{
		http.Port(int(portVal)),
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
