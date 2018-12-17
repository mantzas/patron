package http

import (
	"net/http"
)

// HealthStatus type representing the health of the service via HTTP component.
type HealthStatus int

const (
	// Initializing represents a state when warming up and before the component is Healthy.
	Initializing HealthStatus = 0
	// Healthy represents a state defining a healthy state.
	Healthy HealthStatus = 1
	// Unhealthy represents a state defining a unhealthy state.
	Unhealthy HealthStatus = 2
)

// HealthCheckFunc defines a function type for implementing a health check.
type HealthCheckFunc func() HealthStatus

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
	return NewRouteRaw("/health", http.MethodGet, f, false, nil)
}
