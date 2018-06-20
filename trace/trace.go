package trace

import (
	"io"
	"net/http"
	"time"

	"github.com/mantzas/patron/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
	jaeger "github.com/uber/jaeger-client-go"
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
	tr  opentracing.Tracer
	cls io.Closer
)

// Initialize the tracer if it not already initialized.
func Initialize() {
	if tr != nil {
		return
	}
	tr, cls = jaeger.NewTracer("patron", jaeger.NewConstSampler(true), jaeger.NewNullReporter())
}

// Setup tracing by providing a local agent address.
func Setup(name, agentAddress string) error {
	cfg := config.Configuration{
		ServiceName: name,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentAddress,
		},
	}
	time.Sleep(100 * time.Millisecond)
	metricsFactory := prometheus.New()
	var err error
	tr, cls, err = cfg.NewTracer(
		config.Logger(jaegerLoggerAdapter{}),
		config.Observer(rpcmetrics.NewObserver(metricsFactory.Namespace(name, nil), rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		return errors.Wrap(err, "cannot initialize jaeger tracer")
	}
	return nil
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

// StartConsumerSpan start a new kafka consumer span.
func StartConsumerSpan(name string, cmp Component, hdr map[string]string) opentracing.Span {
	ctx, _ := tr.Extract(opentracing.HTTPHeaders, opentracing.TextMapCarrier(hdr))
	sp := tr.StartSpan(name, consumerOption{ctx: ctx})
	ext.Component.Set(sp, string(cmp))
	return sp
}

// FinishConsumerSpan finished a kafka consumer span.
func FinishConsumerSpan(sp opentracing.Span, hasError bool) {
	ext.Error.Set(sp, hasError)
	sp.Finish()
}

// StartHTTPSpan starts a new HTTP span.
func StartHTTPSpan(path string, r *http.Request) opentracing.Span {
	ctx, _ := tr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	sp := tr.StartSpan(opName(r.Method, path), ext.RPCServerOption(ctx))
	ext.HTTPMethod.Set(sp, r.Method)
	ext.HTTPUrl.Set(sp, r.URL.String())
	ext.Component.Set(sp, "http")
	_ = r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))
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
