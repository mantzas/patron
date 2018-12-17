package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockAuthenticator struct {
	success bool
	err     error
}

func (mo MockAuthenticator) Authenticate(req *http.Request) (bool, error) {
	if mo.err != nil {
		return false, mo.err
	}
	return mo.success, nil
}

func TestNewRoute(t *testing.T) {
	r := NewRoute("/index", http.MethodGet, nil, true, &MockAuthenticator{})
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodGet, r.Method)
	assert.True(t, r.Trace)
	assert.NotNil(t, r.Auth)
}

func TestNewGetRoute(t *testing.T) {
	r := NewGetRoute("/index", nil, true, &MockAuthenticator{})
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodGet, r.Method)
	assert.True(t, r.Trace)
	assert.NotNil(t, r.Auth)
}

func TestNewPostRoute(t *testing.T) {
	r := NewPostRoute("/index", nil, true, &MockAuthenticator{})
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodPost, r.Method)
	assert.True(t, r.Trace)
	assert.NotNil(t, r.Auth)
}

func TestNewPutRoute(t *testing.T) {
	r := NewPutRoute("/index", nil, true, &MockAuthenticator{})
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodPut, r.Method)
	assert.True(t, r.Trace)
	assert.NotNil(t, r.Auth)
}

func TestNewDeleteRoute(t *testing.T) {
	r := NewDeleteRoute("/index", nil, true, &MockAuthenticator{})
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, http.MethodDelete, r.Method)
	assert.True(t, r.Trace)
	assert.NotNil(t, r.Auth)
}
func TestNewRouteRaw(t *testing.T) {
	r := NewRouteRaw("/index", http.MethodGet, nil, false, &MockAuthenticator{})
	assert.Equal(t, "/index", r.Pattern)
	assert.Equal(t, "GET", r.Method)
	assert.False(t, r.Trace)
	assert.NotNil(t, r.Auth)
}
