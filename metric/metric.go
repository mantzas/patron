package metric

import (
	"github.com/thebeatapp/patron/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace = "patron"
)

// Setup metrics.
func Setup(ns string) {
	namespace = ns
}

// NewCounter creates and registers a counter.
func NewCounter(sub, name, help string, labels ...string) (*prometheus.CounterVec, error) {
	cnt := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	if err := prometheus.Register(cnt); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return nil, errors.Wrap(err, "failed to register consumer error metrics")
		}
	}
	return cnt, nil
}

// NewGauge creates and registers a gauge.
func NewGauge(sub, name, help string, labels ...string) (*prometheus.GaugeVec, error) {
	gau := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	if err := prometheus.Register(gau); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return nil, errors.Wrap(err, "failed to register consumer error metrics")
		}
	}
	return gau, nil
}

// NewHistogram creates and registers a histogram.
func NewHistogram(sub, name, help string, labels ...string) (*prometheus.HistogramVec, error) {
	gau := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	if err := prometheus.Register(gau); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return nil, errors.Wrap(err, "failed to register consumer error metrics")
		}
	}
	return gau, nil
}

// NewSummary creates and registers a summary.
func NewSummary(sub, name, help string, labels ...string) (*prometheus.SummaryVec, error) {
	gau := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	if err := prometheus.Register(gau); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			return nil, errors.Wrap(err, "failed to register consumer error metrics")
		}
	}
	return gau, nil
}
