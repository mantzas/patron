package prometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	m := New()
	assert.NotNil(m)
	assert.Len(m.counters, 0)
	assert.Len(m.gauges, 0)
}

func TestMetric_CounterNotExists(t *testing.T) {
	assert := assert.New(t)
	key := "ttt"
	m := New()
	assert.Len(m.counters, 0)
	m.Counter(key, 1.0, "localhost")
}

func TestMetric_RegisterCounterAndAddValue(t *testing.T) {
	assert := assert.New(t)
	key := "test_counter_seconds"
	m := New()
	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: key,
			Help: "Test duration in seconds",
		},
		[]string{"host"},
	)
	m.RegisterCounter(key, c)
	assert.Len(m.counters, 1)
	m.Counter(key, 1.0, "localhost")
}

func TestMetric_GaugeNotExists(t *testing.T) {
	assert := assert.New(t)
	key := "ttt"
	m := New()
	assert.Len(m.gauges, 0)
	m.Gauge(key, 1.0, "localhost")
}

func TestMetric_RegisterGaugeAndAddValue(t *testing.T) {
	assert := assert.New(t)
	key := "test_gauge_seconds"
	m := New()
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: key,
			Help: "Test duration in seconds",
		},
		[]string{"host"},
	)
	m.RegisterGauge(key, g)
	assert.Len(m.gauges, 1)
	m.Gauge(key, 1.0, "localhost")
}

func TestMetric_HistogramNotExists(t *testing.T) {
	assert := assert.New(t)
	key := "ttt"
	m := New()
	assert.Len(m.histograms, 0)
	m.Histogram(key, 1.0, "localhost")
}

func TestMetric_RegisterHistogramAndAddValue(t *testing.T) {
	assert := assert.New(t)
	key := "test_histogram_seconds"
	m := New()
	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: key,
			Help: "Test duration in seconds",
		},
		[]string{"host"},
	)
	m.RegisterHistogram(key, h)
	assert.Len(m.histograms, 1)
	m.Histogram(key, 1.0, "localhost")
}

func TestMetric_SummaryNotExists(t *testing.T) {
	assert := assert.New(t)
	key := "ttt"
	m := New()
	assert.Len(m.summaries, 0)
	m.Summary(key, 1.0, "localhost")
}

func TestMetric_RegisterSummaryAndAddValue(t *testing.T) {
	assert := assert.New(t)
	key := "test_Summary_seconds"
	m := New()
	s := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: key,
			Help: "Test duration in seconds",
		},
		[]string{"host"},
	)
	m.RegisterSummary(key, s)
	assert.Len(m.summaries, 1)
	m.Summary(key, 1.0, "localhost")
}
