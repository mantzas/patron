package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/patron/log"
	patron_http "github.com/mantzas/patron/sync/http"
)

// CreateHandler creates a router.
func CreateHandler(routes []patron_http.Route) http.Handler {

	log.Infof("adding %d routes", len(routes))

	router := httprouter.New()
	for _, route := range routes {
		router.HandlerFunc(route.Method, route.Pattern, route.Handler)
		log.Infof("added route %s %s", route.Method, route.Pattern)
	}
	return router
}

// ParamExtractor extracts parameters from the request.
func ParamExtractor(r *http.Request) map[string]string {
	par := httprouter.ParamsFromContext(r.Context())
	if len(par) == 0 {
		return make(map[string]string, 0)
	}
	p := make(map[string]string, 0)
	for _, v := range par {
		p[v.Key] = v.Value
	}
	return p
}
