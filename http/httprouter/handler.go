package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	patron_http "github.com/mantzas/patron/http"
	"github.com/mantzas/patron/http/middleware"
	"github.com/mantzas/patron/log"
)

// CreateHandler creates a router
func CreateHandler(routes []patron_http.Route) http.Handler {
	log.Infof("adding %d routes", len(routes))

	router := httprouter.New()

	for _, route := range routes {
		h := middleware.DefaultMiddleware(route.Handler)
		router.HandlerFunc(route.Method, route.Pattern, h)
		log.Infof("added route %s %s", route.Method, route.Pattern)
	}

	return router
}
