package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/sync"
	sync_http "github.com/mantzas/patron/sync/http"
)

func process(ctx context.Context, req *sync.Request) (*sync.Response, error) {
	sp, ctx := trace.StartChildSpan(ctx, "google-client", "http-client")
	sp.LogKV("action", "getting www.google.com")
	rsp, err := http.DefaultClient.Get("https://www.google.com")
	if err != nil {
		trace.FinishSpan(sp, true)
		return nil, errors.Wrap(err, "failed to get google.com")
	}
	defer trace.FinishSpan(sp, false)
	return sync.NewResponse(fmt.Sprintf("got %s from google", rsp.Status)), nil
}

func main() {

	// Set up routes
	routes := make([]sync_http.Route, 0)
	routes = append(routes, sync_http.NewRoute("/", http.MethodGet, process, true))

	options := []sync_http.Option{
		sync_http.Port(50000),
		sync_http.Routes(routes),
	}

	httpCp, err := sync_http.New(options...)
	if err != nil {
		log.Fatalf("failed to create HTTP service %v", err)
	}

	srv, err := patron.New("test", httpCp)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}
}
