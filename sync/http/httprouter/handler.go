package httprouter

import (
	"net/http"
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/patron/log"
	patron_http "github.com/mantzas/patron/sync/http"
)

// CreateHandler creates a router
func CreateHandler(routes []patron_http.Route) http.Handler {
	routes = append(routes, profilingRoutes()...)
	log.Infof("adding %d routes", len(routes))

	router := httprouter.New()

	for _, route := range routes {
		router.HandlerFunc(route.Method, route.Pattern, route.Handler)
		log.Infof("added route %s %s", route.Method, route.Pattern)
	}
	return router
}

func profilingRoutes() []patron_http.Route {
	return []patron_http.Route{
		patron_http.NewRouteRaw("/debug/pprof/", http.MethodGet, index),
		patron_http.NewRouteRaw("/debug/pprof/cmdline/", http.MethodGet, cmdline),
		patron_http.NewRouteRaw("/debug/pprof/profile/", http.MethodGet, profile),
		patron_http.NewRouteRaw("/debug/pprof/symbol/", http.MethodGet, symbol),
		patron_http.NewRouteRaw("/debug/pprof/trace/", http.MethodGet, trace),
		patron_http.NewRouteRaw("/debug/pprof/heap/", http.MethodGet, heap),
		patron_http.NewRouteRaw("/debug/pprof/goroutine/", http.MethodGet, goroutine),
		patron_http.NewRouteRaw("/debug/pprof/block/", http.MethodGet, block),
		patron_http.NewRouteRaw("/debug/pprof/threadcreate/", http.MethodGet, threadcreate),
		patron_http.NewRouteRaw("/debug/pprof/mutex/", http.MethodGet, mutex),
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	pprof.Index(w, r)
}

func cmdline(w http.ResponseWriter, r *http.Request) {
	pprof.Cmdline(w, r)
}

func profile(w http.ResponseWriter, r *http.Request) {
	pprof.Profile(w, r)
}

func symbol(w http.ResponseWriter, r *http.Request) {
	pprof.Symbol(w, r)
}

func trace(w http.ResponseWriter, r *http.Request) {
	pprof.Trace(w, r)
}

func heap(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("heap").ServeHTTP(w, r)
}

func goroutine(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("goroutine").ServeHTTP(w, r)
}

func block(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("block").ServeHTTP(w, r)
}

func threadcreate(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("threadcreate").ServeHTTP(w, r)
}

func mutex(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("mutex").ServeHTTP(w, r)
}
