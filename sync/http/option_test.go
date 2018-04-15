package http

import (
	"net/http"
	"testing"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Setup(zerolog.DefaultFactory(log.DebugLevel))
}

func TestSetPorts(t *testing.T) {
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
			err := SetPorts(tt.port)(&s)

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
			err := SetRoutes(tt.rr)(&s)

			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.Equal(tt.rr, s.routes)
				assert.NoError(err)
			}
		})
	}
}
