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
	cls      io.Closer
	innerLog log.Logger
	version  string
)

func init() {
	version = "dev"
}

// Setup tracing by providing all necessary parameters.
func Setup(name, ver, agentAddress, samplerType string, samplerParam float64) error {
	if ver != "" {
		version = ver
	}
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
		config.Logger(jaegerLoggerAdapter{log: log.Create()}),
		config.Observer(rpcmetrics.NewObserver(metricsFactory.Namespace(name, nil), rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		return errors.Wrap(err, "cannot initialize jaeger tracer")
	}
	cls = clsTemp
	opentracing.SetGlobalTracer(tr)
	version = ver
	innerLog = log.Create()
	return nil
}

// Close the tracer.
func Close() error {
	innerLog.Debug("closing tracer")
	return cls.Close()
}

// ConsumerSpan starts a new consumer span.
func ConsumerSpan(
	ctx context.Context,
	name, cmp string,
	hdr map[string]string,
) (opentracing.Span, context.Context) {
	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.TextMapCarrier(hdr))
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		innerLog.Errorf("failed to extract consumer span: %v", err)
	}
	sp := opentracing.StartSpan(name, consumerOption{ctx: spCtx})
	ext.Component.Set(sp, cmp)
	sp.SetTag("version", version)
	return sp, opentracing.ContextWithSpan(ctx, sp)
}

// SpanSuccess finishes a span with a success indicator.
func SpanSuccess(sp opentracing.Span) {
	ext.Error.Set(sp, false)
	sp.Finish()
}

// SpanError finishes a span with a error indicator.
func SpanError(sp opentracing.Span) {
	ext.Error.Set(sp, true)
	sp.Finish()
}

// HTTPSpan starts a new HTTP span.
func HTTPSpan(path string, r *http.Request) (opentracing.Span, *http.Request) {
	ctx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		innerLog.Errorf("failed to extract HTTP span: %v", err)
	}
	sp := opentracing.StartSpan(HTTPOpName(r.Method, path), ext.RPCServerOption(ctx))
	ext.HTTPMethod.Set(sp, r.Method)
	ext.HTTPUrl.Set(sp, r.URL.String())
	ext.Component.Set(sp, "http")
	sp.SetTag("version", version)
	return sp, r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))
}

// FinishHTTPSpan finishes a HTTP span by providing a HTTP status code.
func FinishHTTPSpan(sp opentracing.Span, code int) {
	ext.HTTPStatusCode.Set(sp, uint16(code))
	sp.Finish()
}

// ChildSpan starts a new child span with specified tags.
func ChildSpan(
	ctx context.Context,
	opName, cmp string,
	tags ...opentracing.Tag,
) (opentracing.Span, context.Context) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, opName)
	ext.Component.Set(sp, cmp)
	for _, t := range tags {
		sp.SetTag(t.Key, t.Value)
	}
	sp.SetTag("version", version)
	return sp, ctx
}

// SQLSpan starts a new SQL child span with specified tags.
func SQLSpan(
	ctx context.Context,
	opName, cmp, sqlType, instance, user, stmt string,
	tags ...opentracing.Tag,
) (opentracing.Span, context.Context) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, opName)
	ext.Component.Set(sp, cmp)
	ext.DBType.Set(sp, sqlType)
	ext.DBInstance.Set(sp, instance)
	ext.DBUser.Set(sp, user)
	ext.DBStatement.Set(sp, stmt)
	for _, t := range tags {
		sp.SetTag(t.Key, t.Value)
	}
	sp.SetTag("version", version)
	return sp, ctx
}

// HTTPOpName return a string representation of the HTTP request operation.
func HTTPOpName(method, path string) string {
	return "HTTP " + method + " " + path
}

type jaegerLoggerAdapter struct {
	log log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.log.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.log.Infof(msg, args...)
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
