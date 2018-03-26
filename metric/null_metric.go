package metric

// nullMetric defines a metric that does nothing
type nullMetric struct {
	key    string
	v      float64
	labels []string
}

// Collect does nothing
func (nm *nullMetric) Collect(key string, v float64, labels ...string) {
	nm.key = key
	nm.v = v
	nm.labels = labels
}
