package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/sync"
	patronhttp "github.com/mantzas/patron/sync/http"
	tracehttp "github.com/mantzas/patron/trace/http"
	"github.com/pkg/errors"
)

func init() {
	err := os.Setenv("PATRON_LOG_LEVEL", "debug")
	if err != nil {
		fmt.Printf("failed to set log level env var: %v", err)
		os.Exit(1)
	}
	err = os.Setenv("PATRON_JAEGER_SAMPLER_PARAM", "1.0")
	if err != nil {
		fmt.Printf("failed to set sampler env vars: %v", err)
		os.Exit(1)
	}
}

func main() {
	name := "first"
	version := "1.0.0"

	err := patron.SetupLogging(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	// Set up routes
	routes := []patronhttp.Route{
		patronhttp.NewPostRoute("/", first, true),
	}

	srv, err := patron.New(
		name,
		version,
		patron.Routes(routes),
	)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}
}

func first(ctx context.Context, req *sync.Request) (*sync.Response, error) {
	b, err := json.Encode("patron")
	if err != nil {
		return nil, errors.Wrap(err, "failed create request")
	}

	secondRouteReq, err := http.NewRequest("GET", "http://localhost:50001", bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "failed create request")
	}
	secondRouteReq.Header.Add("Content-Type", "application/json")
	secondRouteReq.Header.Add("Accept", "application/json")
	cl, err := tracehttp.New(tracehttp.Timeout(5 * time.Second))
	if err != nil {
		return nil, err
	}
	rsp, err := cl.Do(ctx, secondRouteReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to post to second service")
	}

	log.Infof("request processed")
	return sync.NewResponse(fmt.Sprintf("got %s from second HTTP route", rsp.Status)), nil
}
