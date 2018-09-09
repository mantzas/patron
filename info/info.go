package info

import (
	"io/ioutil"
	"os"

	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/log"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

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
	Doc     string   `json:"doc,omitempty"`
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

// AddDoc adds documentation from a markdown file.
func AddDoc(file string) {
	if file == "" {
		log.Warn("no file provided")
		serviceInfo.Doc = ""
		return
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Warnf("file %s does not exist", file)
		serviceInfo.Doc = ""
		return
	}
	cnt, err := ioutil.ReadFile(file)
	if err != nil {
		log.Warnf("failed to read file %s", file)
		serviceInfo.Doc = ""
		return
	}
	serviceInfo.Doc = string(blackfriday.Run(cnt, blackfriday.WithExtensions(blackfriday.CommonExtensions)))
}
