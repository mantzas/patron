package http

import (
	"context"
	"net/http"
	"time"

	"github.com/mantzas/patron/trace"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
)

// Client interface of a HTTP client.
type Client interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// TracedClient defines a HTTP client with tracing integrated.
type TracedClient struct {
	cl *http.Client
}

// NewClient creates a new HTTP client.
func NewClient(timeout time.Duration) *TracedClient {
	return &TracedClient{
		cl: &http.Client{
			Timeout:   timeout,
			Transport: &nethttp.Transport{},
		}}
}

// Do executes a HTTP request.
func (tc *TracedClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req, ht := nethttp.TraceRequest(
		opentracing.GlobalTracer(),
		req,
		nethttp.OperationName(trace.HTTPOpName(req.Method, req.URL.String())),
		nethttp.ComponentName(trace.HTTPClientComponent))
	defer ht.Finish()
	return tc.cl.Do(req)
}
