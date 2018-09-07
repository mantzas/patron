package info

// ServiceInfo holds the information of a
var ServiceInfo = Info{}

// Metric describes a metric of the system.
type Metric struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Info contains all information about the service.
type Info struct {
	Name    string   `json:"name,omitempty"`
	Metrics []Metric `json:"metrics,omitempty"`
}

// AddName to the info.
func (i *Info) AddName(n string) {
	i.Name = n
}

// AddMetric to the info.
func (i *Info) AddMetric(m Metric) {
	i.Metrics = append(i.Metrics, m)
}
