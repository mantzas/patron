package route

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	r := New("/index", http.MethodGet, nil)
	assert.Equal("/index", r.Pattern)
	assert.Equal("GET", r.Method)
}
