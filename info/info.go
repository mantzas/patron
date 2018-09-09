package info

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/errors"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// Component information.
type Component struct {
	Type    string            `json:"type,omitempty"`
	Configs map[string]string `json:"configs,omitempty"`
}

// UpsertConfig upsert's the configuration info to the component info.
func (c *Component) UpsertConfig(name, value string) {
	if c.Configs == nil {
		c.Configs = make(map[string]string)
	}
	c.Configs[name] = value
}

type info struct {
	Name       string            `json:"name,omitempty"`
	Version    string            `json:"version,omitempty"`
	Host       string            `json:"host,omitempty"`
	Configs    map[string]string `json:"configs,omitempty"`
	Components []Component       `json:"components,omitempty"`
	Metrics    map[string]string `json:"metrics,omitempty"`
	Doc        string            `json:"doc,omitempty"`
}

// ServiceInfo holds the information of the service.
var serviceInfo = info{
	Configs: make(map[string]string),
	Metrics: make(map[string]string),
}

// Marshal returns the service info as a byte slice.
func Marshal() ([]byte, error) {
	return json.Encode(serviceInfo)
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

// UpsertMetric to the info.
func UpsertMetric(n, d, typ string) {
	serviceInfo.Metrics[n] = fmt.Sprintf("[%s] %s", typ, d)
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

// AppendComponent to the info.
func AppendComponent(cmp Component) {
	serviceInfo.Components = append(serviceInfo.Components, cmp)
}
