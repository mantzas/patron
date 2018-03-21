package metric

import "errors"

// Metric defines a interface that has to be implemented in order to
// be used in this package
type Metric interface {
	Counter(key string, v float64, labels ...string)
	Gauge(key string, v float64, labels ...string)
}

var metric Metric

func init() {
	metric = &nullMetric{}
}

// Setup accepts a implementation of the metric interface
func Setup(m Metric) error {
	if m == nil {
		return errors.New("metric is nil")
	}

	metric = m
	return nil
}

// Counter increases the keyed counter by a value and attaches the labels
func Counter(key string, v float64, labels ...string) {
	metric.Counter(key, v, labels...)
}

// Gauge sets the keyed gauge by a value and attaches the labels
func Gauge(key string, v float64, labels ...string) {
	metric.Gauge(key, v, labels...)
}
