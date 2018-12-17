package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func metricRoute() Route {
	return NewRouteRaw("/metrics", http.MethodGet, promhttp.Handler().ServeHTTP, false)
}
