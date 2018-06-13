package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func metricRoute() Route {
	return NewRouteRaw("/metric", http.MethodGet, promhttp.Handler().ServeHTTP)
}
