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

func TestNewCounterDoubleRegister(t *testing.T) {
	cnt, err := NewCounter("sub", "name", "help", "labels1", "label2")
	assert.NoError(t, err)
	assert.NotNil(t, cnt)
	cnt, err = NewCounter("sub", "name", "help", "labels1", "label2")
	assert.NoError(t, err)
	assert.NotNil(t, cnt)
}

func TestNewGaugeDoubleRegister(t *testing.T) {
	cnt, err := NewGauge("sub", "name", "help", "labels1", "label2")
	assert.NoError(t, err)
	assert.NotNil(t, cnt)
	cnt, err = NewGauge("sub", "name", "help", "labels1", "label2")
	assert.NoError(t, err)
	assert.NotNil(t, cnt)
}

func TestNewHistogramDoubleRegister(t *testing.T) {
	cnt, err := NewHistogram("sub", "name", "help", "labels1", "label2")
	assert.NoError(t, err)
	assert.NotNil(t, cnt)
	cnt, err = NewHistogram("sub", "name", "help", "labels1", "label2")
	assert.NoError(t, err)
	assert.NotNil(t, cnt)
}

func TestNewSummaryDoubleRegister(t *testing.T) {
	cnt, err := NewSummary("sub", "name", "help", "labels1", "label2")
	assert.NoError(t, err)
	assert.NotNil(t, cnt)
	cnt, err = NewSummary("sub", "name", "help", "labels1", "label2")
	assert.NoError(t, err)
	assert.NotNil(t, cnt)
}
