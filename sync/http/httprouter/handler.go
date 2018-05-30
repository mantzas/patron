package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/patron/log"
	patron_http "github.com/mantzas/patron/sync/http"
)

// CreateHandler creates a router
func CreateHandler(routes []patron_http.Route) http.Handler {

	log.Infof("adding %d routes", len(routes))

	router := httprouter.New()
	for _, route := range routes {
		router.HandlerFunc(route.Method, route.Pattern, route.Handler)
		log.Infof("added route %s %s", route.Method, route.Pattern)
	}
	return router
}
