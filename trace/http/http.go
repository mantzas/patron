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

// New creates a new HTTP client.
func New(oo ...OptionFunc) (*TracedClient, error) {
	tc := &TracedClient{
		cl: &http.Client{
			Timeout:   60 * time.Second,
			Transport: &nethttp.Transport{},
		},
	}

	for _, o := range oo {
		err := o(tc)
		if err != nil {
			return nil, err
		}
	}

	return tc, nil
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
