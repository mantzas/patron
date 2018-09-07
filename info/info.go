package info

import "github.com/mantzas/patron/encoding/json"

// ServiceInfo holds the information of the service.
var serviceInfo = info{}

// Marshal returns the service info as a byte slice.
func Marshal() ([]byte, error) {
	return json.Encode(serviceInfo)
}

type metric struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type info struct {
	Name    string   `json:"name,omitempty"`
	Version string   `json:"version,omitempty"`
	Metrics []metric `json:"metrics,omitempty"`
}

// AddName to the info.
func AddName(n string) {
	serviceInfo.Name = n
}

// AddVersion to the info.
func AddVersion(v string) {
	serviceInfo.Version = v
}

// AddMetric to the info.
func AddMetric(n, d string) {
	serviceInfo.Metrics = append(serviceInfo.Metrics, metric{Name: n, Description: d})
}
