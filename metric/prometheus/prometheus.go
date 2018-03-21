package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metric implementation of the metric interface
type Metric struct {
	counters map[string]*prometheus.CounterVec
	gauges   map[string]*prometheus.GaugeVec
}

// New creates a new metric
func New() *Metric {
	return &Metric{
		make(map[string]*prometheus.CounterVec),
		make(map[string]*prometheus.GaugeVec),
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

// Gauge  adds metrics to the keyed gauge and attaches labels
func (pm *Metric) Gauge(key string, v float64, labels ...string) {
	g, ok := pm.gauges[key]
	if !ok {
		return
	}
	g.WithLabelValues(labels...).Set(v)
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
