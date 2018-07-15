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
	synchttp "github.com/mantzas/patron/sync/http"
	"github.com/mantzas/patron/trace/amqp"
	tracehttp "github.com/mantzas/patron/trace/http"
)

type processor struct {
	pub amqp.Publisher
}

func newProcessor(url, exchange string) (*processor, error) {
	p, err := amqp.NewPublisher(url, exchange)
	if err != nil {
		return nil, err
	}
	return &processor{pub: p}, nil
}

func (p *processor) process(ctx context.Context, req *sync.Request) (*sync.Response, error) {
	googleReq, err := http.NewRequest("GET", "https://www.google.com", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create requestfor www.google.com")
	}
	rsp, err := tracehttp.NewClient(5*time.Second).Do(ctx, googleReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get www.google.com")
	}

	msg, err := amqp.NewJSONMessage("test")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create json amqp message")
	}

	err = p.pub.Publish(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to publish amqp message")
	}

	return sync.NewResponse(fmt.Sprintf("got %s from google", rsp.Status)), nil
}

func main() {

	proc, err := newProcessor("amqp://admin:admin@localhost:5672/", "patron")
	if err != nil {
		log.Fatalf("failed to create processor %v", err)
	}

	// Set up routes
	routes := make([]synchttp.Route, 0)
	routes = append(routes, synchttp.NewRoute("/", http.MethodGet, proc.process, true))

	srv, err := patron.New("patron", "1.0.0", patron.Routes(routes))
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}
}
