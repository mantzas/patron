package http

import (
	"context"
	"net/http"
	"time"

	"github.com/mantzas/patron/trace"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

// Do executes a HTTP request with integrated tracing and tracing propagation downstream.
func (tc *TracedClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req, ht := nethttp.TraceRequest(
		opentracing.GlobalTracer(),
		req,
		nethttp.OperationName(trace.HTTPOpName("Client", req.Method, req.URL.String())),
		nethttp.ComponentName(trace.HTTPClientComponent))
	defer ht.Finish()
	rsp, err := tc.cl.Do(req)
	if err != nil {
		ext.Error.Set(ht.Span(), true)
	}
	ext.HTTPMethod.Set(ht.Span(), req.Method)
	ext.HTTPUrl.Set(ht.Span(), req.URL.String())
	ext.HTTPStatusCode.Set(ht.Span(), uint16(rsp.StatusCode))
	return rsp, err
}
