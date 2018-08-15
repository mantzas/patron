package trace

import (
	"context"
	"net/http"
	"testing"

	"github.com/opentracing/opentracing-go"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestSetup_Tracer_Close(t *testing.T) {
	assert := assert.New(t)
	err := Setup("TEST", "1.0.0", "0.0.0.0:6831", "const", 1)
	assert.NoError(err)
	err = Close()
	assert.NoError(err)
	version = "dev"
}

func TestStartFinishConsumerSpan(t *testing.T) {
	assert := assert.New(t)
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	hdr := map[string]string{"key": "val"}
	sp, ctx := ConsumerSpan(context.Background(), "123", AMQPConsumerComponent, hdr)
	assert.NotNil(sp)
	assert.NotNil(ctx)
	assert.IsType(&mocktracer.MockSpan{}, sp)
	jsp := sp.(*mocktracer.MockSpan)
	assert.NotNil(jsp)
	assert.Equal("123", jsp.OperationName)
	SpanError(sp)
	assert.NotNil(sp)
	rawSpan := mtr.FinishedSpans()[0]
	assert.Equal(map[string]interface{}{
		"span.kind": ext.SpanKindConsumerEnum,
		"component": "amqp-consumer",
		"error":     true,
		"version":   "dev",
	}, rawSpan.Tags())
}

func TestStartFinishChildSpan(t *testing.T) {
	assert := assert.New(t)
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	sp, ctx := ConsumerSpan(context.Background(), "123", AMQPConsumerComponent, nil)
	assert.NotNil(sp)
	assert.NotNil(ctx)
	childSp, childCtx := ChildSpan(ctx, "123", "cmp", opentracing.Tag{Key: "key", Value: "value"})
	assert.NotNil(childSp)
	assert.NotNil(childCtx)
	childSp.LogKV("log event")
	assert.IsType(&mocktracer.MockSpan{}, childSp)
	jsp := childSp.(*mocktracer.MockSpan)
	assert.NotNil(jsp)
	assert.Equal("123", jsp.OperationName)
	SpanError(childSp)
	assert.NotNil(childSp)
	rawSpan := mtr.FinishedSpans()[0]
	assert.Equal(map[string]interface{}{
		"component": "cmp",
		"error":     true,
		"key":       "value",
		"version":   "dev",
	}, rawSpan.Tags())
	SpanSuccess(sp)
	rawSpan = mtr.FinishedSpans()[1]
	assert.Equal(map[string]interface{}{
		"component": "amqp-consumer",
		"error":     false,
		"version":   "dev",
		"span.kind": ext.SpanKindConsumerEnum,
	}, rawSpan.Tags())
}

func TestHTTPStartFinishSpan(t *testing.T) {
	assert := assert.New(t)
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	sp, req := HTTPSpan("/", req)
	assert.NotNil(sp)
	assert.NotNil(req)
	assert.IsType(&mocktracer.MockSpan{}, sp)
	jsp := sp.(*mocktracer.MockSpan)
	assert.NotNil(jsp)
	assert.Equal("Server HTTP GET /", jsp.OperationName)
	FinishHTTPSpan(jsp, 200)
	assert.NotNil(jsp)
	rawSpan := mtr.FinishedSpans()[0]
	assert.Equal(map[string]interface{}{
		"span.kind":        ext.SpanKindRPCServerEnum,
		"component":        "http",
		"http.method":      "GET",
		"http.status_code": uint16(200),
		"http.url":         "/",
		"version":          "dev",
	}, rawSpan.Tags())
}
