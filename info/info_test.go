package info

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoMarshal(t *testing.T) {
	expected := `{"metrics":[{"name":"Name","description":"Description"}]}`
	i := Info{}
	i.AddMetric(Metric{Name: "Name", Description: "Description"})
	got, err := json.Marshal(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(got))
}

func TestInfo(t *testing.T) {
	i := Info{}
	i.AddName("Name")
	i.AddMetric(Metric{Name: "Name", Description: "Description"})
	assert.Equal(t, i.Name, "Name")
	assert.Len(t, i.Metrics, 1)
}
