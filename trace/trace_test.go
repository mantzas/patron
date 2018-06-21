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
	err := Setup("TEST", "0.0.0.0:6831", "const", 1)
	assert.NoError(err)
	err = Close()
	assert.NoError(err)
}

func TestStartFinishConsumerSpan(t *testing.T) {
	assert := assert.New(t)
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	hdr := map[string]string{"key": "val"}
	sp := StartConsumerSpan("test", AMQPConsumerComponent, hdr)
	assert.NotNil(sp)
	assert.IsType(&mocktracer.MockSpan{}, sp)
	jsp := sp.(*mocktracer.MockSpan)
	assert.NotNil(jsp)
	assert.Equal("test", jsp.OperationName)
	FinishSpan(sp, true)
	assert.NotNil(sp)
	rawSpan := mtr.FinishedSpans()[0]
	assert.Equal(map[string]interface{}{
		"span.kind": ext.SpanKindConsumerEnum,
		"component": "amqp-consumer",
		"error":     true,
	}, rawSpan.Tags())
}

func TestStartFinishChildSpan(t *testing.T) {
	assert := assert.New(t)
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	sp, ctx := StartChildSpan(context.Background(), "opName", "cmp", opentracing.Tag{Key: "key", Value: "value"})
	assert.NotNil(sp)
	assert.NotNil(ctx)
	sp.LogKV("log event")
	assert.IsType(&mocktracer.MockSpan{}, sp)
	jsp := sp.(*mocktracer.MockSpan)
	assert.NotNil(jsp)
	assert.Equal("opName", jsp.OperationName)
	FinishSpan(sp, true)
	assert.NotNil(sp)
	rawSpan := mtr.FinishedSpans()[0]
	assert.Equal(map[string]interface{}{
		"component": "cmp",
		"error":     true,
		"key":       "value",
	}, rawSpan.Tags())
}

func TestHTTPStartFinishSpan(t *testing.T) {
	assert := assert.New(t)
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	sp, req := StartHTTPSpan("/", req)
	assert.NotNil(sp)
	assert.NotNil(req)
	assert.IsType(&mocktracer.MockSpan{}, sp)
	jsp := sp.(*mocktracer.MockSpan)
	assert.NotNil(jsp)
	assert.Equal("HTTP GET /", jsp.OperationName)
	FinishHTTPSpan(jsp, 200)
	assert.NotNil(jsp)
	rawSpan := mtr.FinishedSpans()[0]
	assert.Equal(map[string]interface{}{
		"span.kind":        ext.SpanKindRPCServerEnum,
		"component":        "http",
		"http.method":      "GET",
		"http.status_code": uint16(200),
		"http.url":         "/",
	}, rawSpan.Tags())
}
