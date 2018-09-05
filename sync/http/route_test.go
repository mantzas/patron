package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoute(t *testing.T) {
	assert := assert.New(t)
	r := NewRoute("/index", http.MethodGet, nil, true)
	assert.Equal("/index", r.Pattern)
	assert.Equal(http.MethodGet, r.Method)
	assert.True(r.Trace)
}

func TestNewGetRoute(t *testing.T) {
	assert := assert.New(t)
	r := NewGetRoute("/index", nil, true)
	assert.Equal("/index", r.Pattern)
	assert.Equal(http.MethodGet, r.Method)
	assert.True(r.Trace)
}

func TestNewPostRoute(t *testing.T) {
	assert := assert.New(t)
	r := NewPostRoute("/index", nil, true)
	assert.Equal("/index", r.Pattern)
	assert.Equal(http.MethodPost, r.Method)
	assert.True(r.Trace)
}

func TestNewPutRoute(t *testing.T) {
	assert := assert.New(t)
	r := NewPutRoute("/index", nil, true)
	assert.Equal("/index", r.Pattern)
	assert.Equal(http.MethodPut, r.Method)
	assert.True(r.Trace)
}

func TestNewDeleteRoute(t *testing.T) {
	assert := assert.New(t)
	r := NewDeleteRoute("/index", nil, true)
	assert.Equal("/index", r.Pattern)
	assert.Equal(http.MethodDelete, r.Method)
	assert.True(r.Trace)
}
func TestNewRouteRaw(t *testing.T) {
	assert := assert.New(t)
	r := NewRouteRaw("/index", http.MethodGet, nil)
	assert.Equal("/index", r.Pattern)
	assert.Equal("GET", r.Method)
	assert.False(r.Trace)
}
