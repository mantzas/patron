package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/sync"
	sync_http "github.com/mantzas/patron/sync/http"
	trace_http "github.com/mantzas/patron/trace/http"
)

func process(ctx context.Context, req *sync.Request) (*sync.Response, error) {
	googleReq, err := http.NewRequest("GET", "https://www.google.com", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create requestfor www.google.com")
	}
	rsp, err := trace_http.NewClient(1*time.Second).Do(ctx, googleReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get www.google.com")
	}
	return sync.NewResponse(fmt.Sprintf("got %s from google", rsp.Status)), nil
}

func main() {

	// Set up routes
	routes := make([]sync_http.Route, 0)
	routes = append(routes, sync_http.NewRoute("/", http.MethodGet, process, true))

	srv, err := patron.New("test", patron.Routes(routes))
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}
}
