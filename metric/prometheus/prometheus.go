package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metric implementation of the metric interface
type Metric struct {
	counters   map[string]*prometheus.CounterVec
	gauges     map[string]*prometheus.GaugeVec
	histograms map[string]*prometheus.HistogramVec
	summaries  map[string]*prometheus.SummaryVec
}

// New creates a new metric
func New() *Metric {
	return &Metric{
		make(map[string]*prometheus.CounterVec),
		make(map[string]*prometheus.GaugeVec),
		make(map[string]*prometheus.HistogramVec),
		make(map[string]*prometheus.SummaryVec),
	}
}

// Counter adds metrics to the keyed counter and attaches labels
func (pm *Metric) Counter(key string, v float64, labels ...string) {
	c, ok := pm.counters[key]
	if !ok {
		return
	}
	c.WithLabelValues(labels...).Add(v)
}

// Gauge adds metrics to the keyed gauge and attaches labels
func (pm *Metric) Gauge(key string, v float64, labels ...string) {
	g, ok := pm.gauges[key]
	if !ok {
		return
	}
	g.WithLabelValues(labels...).Set(v)
}

// Histogram adds metrics to the keyed histogram and attaches labels
func (pm *Metric) Histogram(key string, v float64, labels ...string) {
	h, ok := pm.histograms[key]
	if !ok {
		return
	}
	h.WithLabelValues(labels...).Observe(v)
}

// Summary adds metrics to the keyed summary and attaches labels
func (pm *Metric) Summary(key string, v float64, labels ...string) {
	s, ok := pm.summaries[key]
	if !ok {
		return
	}
	s.WithLabelValues(labels...).Observe(v)
}

// RegisterCounter registers a counter
func (pm *Metric) RegisterCounter(key string, c *prometheus.CounterVec) {
	prometheus.MustRegister(c)
	pm.counters[key] = c
}

// RegisterGauge registers a gauge
func (pm *Metric) RegisterGauge(key string, g *prometheus.GaugeVec) {
	prometheus.MustRegister(g)
	pm.gauges[key] = g
}

// RegisterHistogram registers a histogram
func (pm *Metric) RegisterHistogram(key string, h *prometheus.HistogramVec) {
	prometheus.MustRegister(h)
	pm.histograms[key] = h
}

// RegisterSummary registers a summary
func (pm *Metric) RegisterSummary(key string, s *prometheus.SummaryVec) {
	prometheus.MustRegister(s)
	pm.summaries[key] = s
}
