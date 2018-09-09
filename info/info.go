package info

import (
	"io/ioutil"
	"os"

	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/errors"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// ServiceInfo holds the information of the service.
var serviceInfo = info{
	Configs: make(map[string]string),
	Metrics: make(map[string]string),
}

// Marshal returns the service info as a byte slice.
func Marshal() ([]byte, error) {
	return json.Encode(serviceInfo)
}

type info struct {
	Name    string            `json:"name,omitempty"`
	Version string            `json:"version,omitempty"`
	Host    string            `json:"host,omitempty"`
	Metrics map[string]string `json:"metrics,omitempty"`
	Configs map[string]string `json:"configs,omitempty"`
	Doc     string            `json:"doc,omitempty"`
}

// UpdateName to the info.
func UpdateName(n string) {
	serviceInfo.Name = n
}

// UpdateVersion to the info.
func UpdateVersion(v string) {
	serviceInfo.Version = v
}

// UpdateHost to the info.
func UpdateHost(h string) {
	serviceInfo.Host = h
}

// AddMetric to the info.
func AddMetric(n, d string) {
	serviceInfo.Metrics[n] = d
}

// ImportDoc adds documentation from a markdown file.
func ImportDoc(file string) error {
	if file == "" {
		serviceInfo.Doc = ""
		return errors.New("no file provided")
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		serviceInfo.Doc = ""
		return errors.Errorf("file %s does not exist", file)
	}
	cnt, err := ioutil.ReadFile(file)
	if err != nil {
		serviceInfo.Doc = ""
		return errors.Errorf("failed to read file %s", file)
	}
	serviceInfo.Doc = string(blackfriday.Run(cnt, blackfriday.WithExtensions(blackfriday.CommonExtensions)))
	return nil
}

// UpsertConfig to the info.
func UpsertConfig(n, v string) {
	serviceInfo.Configs[n] = v
}
