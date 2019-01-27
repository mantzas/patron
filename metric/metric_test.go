package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	expected := "namespace"
	Setup(expected)
	assert.Equal(t, expected, namespace)
}

func TestNewCounterDoubleRegisterPanics(t *testing.T) {
	m1 := NewCounter("sub", "name", "help", "labels1", "label2")
	m2 := NewCounter("sub", "name", "help", "labels1", "label2")
	assert.Panics(t, func() {
		MustRegister(m1, m2)
	})
}

func TestNewGaugeDoubleRegisterPanics(t *testing.T) {
	m1 := NewGauge("sub", "name", "help", "labels1", "label2")
	m2 := NewGauge("sub", "name", "help", "labels1", "label2")
	assert.Panics(t, func() {
		MustRegister(m1, m2)
	})
}

func TestNewHistogramDoubleRegisterPanics(t *testing.T) {
	m1 := NewHistogram("sub", "name", "help", "labels1", "label2")
	m2 := NewHistogram("sub", "name", "help", "labels1", "label2")
	assert.Panics(t, func() {
		MustRegister(m1, m2)
	})
}

func TestNewSummaryDoubleRegisterPanics(t *testing.T) {
	m1 := NewSummary("sub", "name", "help", "labels1", "label2")
	m2 := NewSummary("sub", "name", "help", "labels1", "label2")
	assert.Panics(t, func() {
		MustRegister(m1, m2)
	})
}
