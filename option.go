package patron

import (
	"time"

	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

const (
	minReportingPeriod = 1 * time.Second
)

// Option defines a option for the HTTP service
type Option func(*Server) error

// Metric option for setting up metrics with exporter and reporting duration
func Metric(e view.Exporter, rp time.Duration) Option {
	return func(s *Server) error {
		if e == nil {
			return errors.New("exporter is not defined")
		}
		if rp < minReportingPeriod {
			return errors.New("reporting period is too small")
		}
		view.RegisterExporter(e)
		view.SetReportingPeriod(rp)
		log.Infof("metric set with report duration %v", rp)
		return nil
	}
}

// Trace option for setting up metrics with exporter and config
func Trace(e trace.Exporter, cfg trace.Config) Option {
	return func(s *Server) error {
		if e == nil {
			return errors.New("exporter is not defined")
		}
		trace.RegisterExporter(e)
		trace.ApplyConfig(cfg)
		log.Info("trace set")
		return nil
	}
}
