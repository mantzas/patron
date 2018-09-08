package info

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	AddName("Name")
	AddVersion("1.2.3")
	AddMetric("Name", "Description")
	expected := `{"name":"Name","version":"1.2.3","metrics":[{"name":"Name","description":"Description"`
	got, err := Marshal()
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(got), expected))
}
