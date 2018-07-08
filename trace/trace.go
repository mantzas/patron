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

const (
	// KafkaConsumerComponent definition.
	KafkaConsumerComponent = "kafka-consumer"
	// KafkaAsyncProducerComponent definition.
	KafkaAsyncProducerComponent = "kafka-async-producer"
	// AMQPConsumerComponent definition.
	AMQPConsumerComponent = "amqp-consumer"
	// AMQPPublisherComponent definition.
	AMQPPublisherComponent = "amqp-publisher"
	// HTTPComponent definition.
	HTTPComponent = "http"
	// HTTPClientComponent definition.
	HTTPClientComponent = "http-client"
)

var (
	cls io.Closer
)

// Setup tracing by providing all necessary parameters.
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

// StartConsumerSpan starts a new consumer span.
func StartConsumerSpan(ctx context.Context, name, cmp string, hdr map[string]string) (opentracing.Span, context.Context) {
	spCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.TextMapCarrier(hdr))
	sp := opentracing.StartSpan(name, consumerOption{ctx: spCtx})
	ext.Component.Set(sp, cmp)
	return sp, opentracing.ContextWithSpan(ctx, sp)
}

// FinishSpanWithSuccess finishes a span with a success indicator.
func FinishSpanWithSuccess(sp opentracing.Span) {
	ext.Error.Set(sp, false)
	sp.Finish()
}

// FinishSpanWithError finishes a span with a error indicator.
func FinishSpanWithError(sp opentracing.Span) {
	ext.Error.Set(sp, true)
	sp.Finish()
}

// StartHTTPSpan starts a new HTTP span.
func StartHTTPSpan(path string, r *http.Request) (opentracing.Span, *http.Request) {
	ctx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	sp := opentracing.StartSpan(HTTPOpName(r.Method, path), ext.RPCServerOption(ctx))
	ext.HTTPMethod.Set(sp, r.Method)
	ext.HTTPUrl.Set(sp, r.URL.String())
	ext.Component.Set(sp, "http")
	return sp, r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))
}

// FinishHTTPSpan finishes a HTTP span by providing a HTTP status code.
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

// HTTPOpName return a string representation of the HTTP request operation.
func HTTPOpName(method, path string) string {
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
