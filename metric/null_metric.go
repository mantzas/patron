package metric

// nullMetric defines a metric that does nothing
type nullMetric struct {
	key    string
	v      float64
	labels []string
}

// Counter does nothing
func (nm *nullMetric) Counter(key string, v float64, labels ...string) {
	nm.key = key
	nm.v = v
	nm.labels = labels
}

// Gauge does nothing
func (nm *nullMetric) Gauge(key string, v float64, labels ...string) {
	nm.key = key
	nm.v = v
	nm.labels = labels
}
