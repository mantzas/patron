package info

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	i := info{}
	i.AddName("Name")
	i.AddMetric("Name", "Description")
	assert.Equal(t, i.Name, "Name")
	assert.Len(t, i.Metrics, 1)
	expected := `{"name":"Name","metrics":[{"name":"Name","description":"Description"}]}`
	got, err := json.Marshal(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(got))
}
