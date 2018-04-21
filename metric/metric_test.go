package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		f       Metric
		wantErr bool
	}{
		{"failure with nil metric", nil, true},
		{"success", &nullMetric{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Setup(tt.f)
			if tt.wantErr {
				assert.Error(err, "expected error")
			} else {
				assert.NoError(err, "error not expected")
			}
		})
	}
}

func TestCollect(t *testing.T) {
	assert := assert.New(t)
	m := nullMetric{}
	Setup(&m)
	key := "key"
	value := 1.99
	labels := []string{"test 1", "test 2"}
	Collect(key, value, labels...)
	assert.Equal(key, m.key)
	assert.Equal(value, m.v)
	assert.Equal(labels, m.labels)
}
