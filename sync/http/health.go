package http

import (
	"net/http"
)

// HealthStatus type represanting the health of a service
type HealthStatus int

const (
	// Initializing represents a state before the service is Healthy
	Initializing HealthStatus = 0
	// Healthy represents a state defining a healthy state
	Healthy HealthStatus = 1
	// Unhealthy represents a state defining a unhealthy state
	Unhealthy HealthStatus = 2
)

// HealthCheckFunc defines a function for implementing a health check
type HealthCheckFunc func() HealthStatus

// HealthCheckRoute returns a route for implementing a health check
func healthCheckRoute(hcf HealthCheckFunc) Route {

	f := func(w http.ResponseWriter, r *http.Request) {
		switch hcf() {
		case Initializing:
			w.WriteHeader(http.StatusServiceUnavailable)
		case Healthy:
			w.WriteHeader(http.StatusOK)
		case Unhealthy:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}

	return NewRouteRaw("/health", http.MethodGet, f)
}
