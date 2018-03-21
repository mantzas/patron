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
