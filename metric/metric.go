package metric

import "errors"

// Metric defines a interface that has to be implemented in order to
// be used in this package
type Metric interface {
	Collect(key string, v float64, labels ...string)
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

// Collect a keyed value with labels
func Collect(key string, v float64, labels ...string) {
	metric.Collect(key, v, labels...)
}
