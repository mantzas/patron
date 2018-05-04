package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPort(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"success", 30000, false},
		{"error for port number out of range", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{}
			err := Port(tt.port)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.Equal(tt.port, s.port)
				assert.NoError(err)
			}
		})
	}
}

func TestSetRoutes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		rr      []Route
		wantErr bool
	}{
		{"success", []Route{NewRoute("/", http.MethodGet, nil)}, false},
		{"error for no routes", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{}
			err := Routes(tt.rr)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.Equal(tt.rr, s.routes)
				assert.NoError(err)
			}
		})
	}
}

func TestSetHealthCheck(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		hcf     HealthCheckFunc
		wantErr bool
	}{
		{"success", func() HealthStatus { return Healthy }, false},
		{"error for no routes", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{}
			err := HealthCheck(tt.hcf)(&s)
			if tt.wantErr {
				assert.Error(err)
				assert.Nil(s.hc)
			} else {
				assert.NoError(err)
				assert.NotNil(s.hc)
			}
		})
	}
}
