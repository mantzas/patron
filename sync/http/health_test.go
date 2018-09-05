package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_healthCheckRoute(t *testing.T) {
	tests := []struct {
		name string
		hcf  HealthCheckFunc
		want int
	}{
		{"healthy", func() HealthStatus { return Healthy }, http.StatusOK},
		{"initializing", func() HealthStatus { return Initializing }, http.StatusServiceUnavailable},
		{"unhealthy", func() HealthStatus { return Unhealthy }, http.StatusInternalServerError},
		{"default", func() HealthStatus { return 10 }, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := healthCheckRoute(tt.hcf)
			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/health", nil)
			assert.NoError(t, err)
			r.Handler(resp, req)
			assert.Equal(t, tt.want, resp.Code)
		})
	}
}
