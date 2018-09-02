package http

import (
	"testing"
	"time"

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
			s := Component{}
			err := Port(tt.port)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.Equal(tt.port, s.httpPort)
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
		{"success", []Route{NewGetRoute("/", testHandler{}.Process, true)}, false},
		{"error for no routes", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Component{}
			err := Routes(tt.rr)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.Len(s.routes, 1)
				assert.Equal(tt.rr[0].Method, s.routes[0].Method)
				assert.Equal(tt.rr[0].Pattern, s.routes[0].Pattern)
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
			s := Component{}
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

func TestSetSecure(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		certFile string
		keyFile  string
		wantErr  bool
	}{
		{"success", "certFile", "keyFile", false},
		{"failed missing cert file", "", "keyFile", true},
		{"failed missing key file", "certFile", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Component{}
			err := Secure(tt.certFile, tt.keyFile)(&s)
			if tt.wantErr {
				assert.Error(err)
				assert.Empty(s.certFile)
				assert.Empty(s.keyFile)
			} else {
				assert.NoError(err)
				assert.NotEmpty(s.certFile)
				assert.NotEmpty(s.keyFile)
			}
		})
	}
}

func TestTimeouts(t *testing.T) {
	assert := assert.New(t)
	c := Component{}
	err := Timeouts(2*time.Second, 3*time.Second)(&c)
	assert.NoError(err)
	assert.Equal(2*time.Second, c.httpReadTimeout)
	assert.Equal(3*time.Second, c.httpWriteTimeout)
}
