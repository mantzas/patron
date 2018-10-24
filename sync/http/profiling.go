package http

import (
	"net/http"
	"net/http/pprof"
)

func profilingRoutes() []Route {
	return []Route{
		NewRouteRaw("/debug/pprof/", http.MethodGet, profIndex),
		NewRouteRaw("/debug/pprof/allocs/", http.MethodGet, pprofAllocsIndex),
		NewRouteRaw("/debug/pprof/cmdline/", http.MethodGet, profCmdline),
		NewRouteRaw("/debug/pprof/profile/", http.MethodGet, profProfile),
		NewRouteRaw("/debug/pprof/symbol/", http.MethodGet, profSymbol),
		NewRouteRaw("/debug/pprof/trace/", http.MethodGet, profTrace),
		NewRouteRaw("/debug/pprof/heap/", http.MethodGet, profHeap),
		NewRouteRaw("/debug/pprof/goroutine/", http.MethodGet, profGoroutine),
		NewRouteRaw("/debug/pprof/block/", http.MethodGet, profBlock),
		NewRouteRaw("/debug/pprof/threadcreate/", http.MethodGet, profThreadcreate),
		NewRouteRaw("/debug/pprof/mutex/", http.MethodGet, profMutex),
	}
}

func profIndex(w http.ResponseWriter, r *http.Request) {
	pprof.Index(w, r)
}

func pprofAllocsIndex(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("allocs").ServeHTTP(w, r)
}

func profCmdline(w http.ResponseWriter, r *http.Request) {
	pprof.Cmdline(w, r)
}

func profProfile(w http.ResponseWriter, r *http.Request) {
	pprof.Profile(w, r)
}

func profSymbol(w http.ResponseWriter, r *http.Request) {
	pprof.Symbol(w, r)
}

func profTrace(w http.ResponseWriter, r *http.Request) {
	pprof.Trace(w, r)
}

func profHeap(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("heap").ServeHTTP(w, r)
}

func profGoroutine(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("goroutine").ServeHTTP(w, r)
}

func profBlock(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("block").ServeHTTP(w, r)
}

func profThreadcreate(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("threadcreate").ServeHTTP(w, r)
}

func profMutex(w http.ResponseWriter, r *http.Request) {
	pprof.Handler("mutex").ServeHTTP(w, r)
}
