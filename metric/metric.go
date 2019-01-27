package metric

import (
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
func NewCounter(sub, name, help string, labels ...string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		labels,
	)
}

// NewGauge creates and registers a gauge.
func NewGauge(sub, name, help string, labels ...string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		labels,
	)
}

// NewHistogram creates and registers a histogram.
func NewHistogram(sub, name, help string, labels ...string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		labels,
	)
}

// NewSummary creates and registers a summary.
func NewSummary(sub, name, help string, labels ...string) *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: sub,
			Name:      name,
			Help:      help,
		},
		labels,
	)
}

// MustRegister registers the collector, panics on error.
func MustRegister(cs ...prometheus.Collector) {
	prometheus.MustRegister(cs...)
}
