package info

// ServiceInfo holds the information of the service.
var ServiceInfo = info{}

type metric struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type info struct {
	Name    string   `json:"name,omitempty"`
	Metrics []metric `json:"metrics,omitempty"`
}

// AddName to the info.
func (i *info) AddName(n string) {
	i.Name = n
}

// AddMetric to the info.
func (i *info) AddMetric(n, d string) {
	i.Metrics = append(i.Metrics, metric{Name: n, Description: d})
}
