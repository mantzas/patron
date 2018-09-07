package info

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	AddName("Name")
	AddMetric("Name", "Description")
	expected := `{"name":"Name","metrics":[{"name":"Name","description":"Description"}]}`
	got, err := Marshal()
	assert.NoError(t, err)
	assert.Equal(t, expected, string(got))
}
