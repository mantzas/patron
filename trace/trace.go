package trace

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/mantzas/patron/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

// Component enum definition.
type Component string

const (
	// KafkaConsumerComponent definition.
	KafkaConsumerComponent Component = "kafka-consumer"
	// AMQPConsumerComponent definition.
	AMQPConsumerComponent Component = "amqp-consumer"
	// HTTPComponent definition.
	HTTPComponent Component = "http"
)

var (
	cls io.Closer
)

// Setup tracing by providing a local agent address.
func Setup(name, agentAddress, samplerType string, samplerParam float64) error {
	cfg := config.Configuration{
		ServiceName: name,
		Sampler: &config.SamplerConfig{
			Type:  samplerType,
			Param: samplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentAddress,
		},
	}
	time.Sleep(100 * time.Millisecond)
	metricsFactory := prometheus.New()
	tr, clsTemp, err := cfg.NewTracer(
		config.Logger(jaegerLoggerAdapter{}),
		config.Observer(rpcmetrics.NewObserver(metricsFactory.Namespace(name, nil), rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		return errors.Wrap(err, "cannot initialize jaeger tracer")
	}
	cls = clsTemp
	opentracing.SetGlobalTracer(tr)
	return nil
}

// Close the tracer.
func Close() error {
	log.Info("closing tracer")
	return cls.Close()
}

// StartConsumerSpan start a new kafka consumer span.
func StartConsumerSpan(ctx context.Context, name string, cmp Component, hdr map[string]string) (opentracing.Span, context.Context) {
	spCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.TextMapCarrier(hdr))
	sp := opentracing.StartSpan(name, consumerOption{ctx: spCtx})
	ext.Component.Set(sp, string(cmp))
	return sp, opentracing.ContextWithSpan(ctx, sp)
}

// FinishSpan finished a kafka consumer span.
func FinishSpan(sp opentracing.Span, hasError bool) {
	ext.Error.Set(sp, hasError)
	sp.Finish()
}

// StartHTTPSpan starts a new HTTP span.
func StartHTTPSpan(path string, r *http.Request) (opentracing.Span, *http.Request) {
	ctx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	sp := opentracing.StartSpan(opName(r.Method, path), ext.RPCServerOption(ctx))
	ext.HTTPMethod.Set(sp, r.Method)
	ext.HTTPUrl.Set(sp, r.URL.String())
	ext.Component.Set(sp, "http")
	return sp, r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))
}

// FinishHTTPSpan finishes a HTTP span.
func FinishHTTPSpan(sp opentracing.Span, code int) {
	ext.HTTPStatusCode.Set(sp, uint16(code))
	sp.Finish()
}

// StartChildSpan starts a new child span with specified tags.
func StartChildSpan(ctx context.Context, opName, cmp string, tags ...opentracing.Tag) (opentracing.Span, context.Context) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, opName)
	ext.Component.Set(sp, cmp)
	for _, t := range tags {
		sp.SetTag(t.Key, t.Value)
	}

	return sp, ctx
}

func opName(method, path string) string {
	return "HTTP " + method + " " + path
}

type jaegerLoggerAdapter struct {
}

func (l jaegerLoggerAdapter) Error(msg string) {
	log.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

type consumerOption struct {
	ctx opentracing.SpanContext
}

func (r consumerOption) Apply(o *opentracing.StartSpanOptions) {
	if r.ctx != nil {
		opentracing.ChildOf(r.ctx).Apply(o)
	}
	ext.SpanKindConsumer.Apply(o)
}
