package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoute(t *testing.T) {
	assert := assert.New(t)
	r := NewRoute("/index", http.MethodGet, nil, nil, true)
	assert.Equal("/index", r.Pattern)
	assert.Equal("GET", r.Method)
	assert.True(r.Trace)
}

func TestNewRouteRaw(t *testing.T) {
	assert := assert.New(t)
	r := NewRouteRaw("/index", http.MethodGet, nil)
	assert.Equal("/index", r.Pattern)
	assert.Equal("GET", r.Method)
	assert.False(r.Trace)
}
