package http

import (
	"net/http"
	"net/http/pprof"
)

func profilingRoutes() []Route {
	return []Route{
		NewRouteRaw("/debug/pprof/", http.MethodGet, index, false),
		NewRouteRaw("/debug/pprof/cmdline/", http.MethodGet, cmdline, false),
		NewRouteRaw("/debug/pprof/profile/", http.MethodGet, profile, false),
		NewRouteRaw("/debug/pprof/symbol/", http.MethodGet, symbol, false),
		NewRouteRaw("/debug/pprof/trace/", http.MethodGet, trace, false),
		NewRouteRaw("/debug/pprof/heap/", http.MethodGet, heap, false),
		NewRouteRaw("/debug/pprof/goroutine/", http.MethodGet, goroutine, false),
		NewRouteRaw("/debug/pprof/block/", http.MethodGet, block, false),
		NewRouteRaw("/debug/pprof/threadcreate/", http.MethodGet, threadcreate, false),
		NewRouteRaw("/debug/pprof/mutex/", http.MethodGet, mutex, false),
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
