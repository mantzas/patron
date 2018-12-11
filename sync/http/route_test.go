package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoute(t *testing.T) {
	r := NewRoute("/index", http.MethodGet, nil, true)
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodGet, r.Method)
	assert.True(t, r.Trace)
}

func TestNewGetRoute(t *testing.T) {
	r := NewGetRoute("/index", nil, true)
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodGet, r.Method)
	assert.True(t, r.Trace)
}

func TestNewPostRoute(t *testing.T) {
	r := NewPostRoute("/index", nil, true)
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodPost, r.Method)
	assert.True(t, r.Trace)
}

func TestNewPutRoute(t *testing.T) {
	r := NewPutRoute("/index", nil, true)
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodPut, r.Method)
	assert.True(t, r.Trace)
}

func TestNewDeleteRoute(t *testing.T) {
	r := NewDeleteRoute("/index", nil, true)
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodDelete, r.Method)
	assert.True(t, r.Trace)
}
func TestNewRouteRaw(t *testing.T) {
	r := NewRouteRaw("/index", http.MethodGet, nil, false)
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, "GET", r.Method)
	assert.False(t, r.Trace)
}
