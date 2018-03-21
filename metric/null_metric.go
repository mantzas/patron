package metric

// nullMetric defines a metric that does nothing
type nullMetric struct {
}

// Counter does nothing
func (nm nullMetric) Counter(key string, v float64, labels ...string) {
}

// Gauge does nothing
func (nm nullMetric) Gauge(key string, v float64, labels ...string) {
}
