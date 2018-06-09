package trace

import (
	"io"
	"net/http"

	"github.com/mantzas/patron/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

var (
	tr  opentracing.Tracer
	cls io.Closer
)

func init() {
	tr, cls = jaeger.NewTracer("patron", jaeger.NewConstSampler(true), jaeger.NewNullReporter())
}

// Setup a new tracer.
func Setup(name string, sampler jaeger.Sampler, reporter jaeger.Reporter, options ...jaeger.TracerOption) {
	log.Info("setting up tracer")
	tr, cls = jaeger.NewTracer(name, sampler, reporter, options...)
}

// Tracer returns the setup tracer.
func Tracer() opentracing.Tracer {
	return tr
}

// Close the tracer.
func Close() error {
	log.Info("closing tracer")
	return cls.Close()
}

// StartHTTPSpan starts a new HTTP span.
func StartHTTPSpan(path string, r *http.Request) opentracing.Span {
	ctx, _ := tr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	sp := tr.StartSpan(opName(r.Method, path), ext.RPCServerOption(ctx))
	ext.HTTPMethod.Set(sp, r.Method)
	ext.HTTPUrl.Set(sp, r.URL.String())
	ext.Component.Set(sp, "http")
	r = r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))
	return sp
}

// FinishHTTPSpan finishes a HTTP span.
func FinishHTTPSpan(sp opentracing.Span, code int) {
	ext.HTTPStatusCode.Set(sp, uint16(code))
	sp.Finish()
}

func opName(method, path string) string {
	return "HTTP " + method + " " + path
}
