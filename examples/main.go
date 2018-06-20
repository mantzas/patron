package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/mantzas/patron/config/env"
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/config"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/mantzas/patron/sync"
	sync_http "github.com/mantzas/patron/sync/http"
)

type indexProcessor struct {
}

func (ih indexProcessor) Process(ctx context.Context, req *sync.Request) (*sync.Response, error) {
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

func init() {

	cfg, err := env.New(nil)
	if err != nil {
		fmt.Printf("failed to setup env config %v", err)
		os.Exit(1)
	}

	// Set up config (should come from flag, env, file etc)
	err = config.Setup(cfg)
	if err != nil {
		fmt.Printf("failed to setup config %v", err)
		os.Exit(1)
	}

	err = config.Set("LOG_LEVEL", "info")
	if err != nil {
		fmt.Printf("failed to set log level config %v", err)
		os.Exit(1)
	}

	err = config.Set("JAEGER_LOCAL_ADDR", "0.0.0.0:6831")
	if err != nil {
		fmt.Printf("failed to set jaeger local address %v", err)
		os.Exit(1)
	}
}

func main() {

	// Set up logging
	lvl, err := config.GetString("LOG_LEVEL")
	if err != nil {
		fmt.Printf("failed to get log level config %v", err)
		os.Exit(1)
	}

	err = log.Setup(zerolog.DefaultFactory(log.Level(lvl)))
	if err != nil {
		fmt.Printf("failed to setup logging %v", err)
		os.Exit(1)
	}

	jaegerAddr, err := config.GetString("JAEGER_LOCAL_ADDR")
	if err != nil {
		log.Fatalf("failed to get jaeger local address %v", err)
	}

	// Set up routes
	routes := make([]sync_http.Route, 0)
	routes = append(routes, sync_http.NewRoute("/", http.MethodGet, indexProcessor{}.Process, true))

	options := []sync_http.Option{
		sync_http.Port(50000),
		sync_http.Routes(routes),
	}

	httpCp, err := sync_http.New(options...)
	if err != nil {
		log.Fatalf("failed to create HTTP service %v", err)
	}

	srv, err := patron.New("test", []patron.Component{httpCp}, patron.Tracing(jaegerAddr))
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}
}
