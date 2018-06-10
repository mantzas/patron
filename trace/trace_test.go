package trace

import (
	"net/http"
	"testing"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestInitialize_Tracer_Close(t *testing.T) {
	assert := assert.New(t)
	Initialize()
	tr := Tracer()
	assert.NotNil(tr)
	Initialize()
	tr1 := Tracer()
	assert.Equal(tr, tr1)
	err := Close()
	assert.NoError(err)
}

func TestSetup_Tracer_Close(t *testing.T) {
	assert := assert.New(t)
	err := Setup("TEST", "0.0.0.0:6831")
	assert.NoError(err)
	tr := Tracer()
	assert.NotNil(tr)
	err = Close()
	assert.NoError(err)
}

func TestStartFinishSpan(t *testing.T) {
	assert := assert.New(t)
	mtr := mocktracer.New()
	tr = mtr
	hdr := map[string]string{"key": "val"}
	sp := StartConsumerSpan("test", AMQPConsumerComponent, hdr)
	assert.NotNil(sp)
	assert.IsType(&mocktracer.MockSpan{}, sp)
	jsp := sp.(*mocktracer.MockSpan)
	assert.NotNil(jsp)
	assert.Equal("test", jsp.OperationName)
	FinishConsumerSpan(sp, true)
	assert.NotNil(sp)
	rawSpan := mtr.FinishedSpans()[0]
	assert.Equal(map[string]interface{}{
		"span.kind": ext.SpanKindConsumerEnum,
		"component": "amqp-consumer",
		"error":     true,
	}, rawSpan.Tags())
}

func TestHTTPStartFinishSpan(t *testing.T) {
	assert := assert.New(t)
	mtr := mocktracer.New()
	tr = mtr
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	sp := StartHTTPSpan("/", req)
	assert.NotNil(sp)
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
