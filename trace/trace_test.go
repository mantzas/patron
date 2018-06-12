package trace

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go"
)

func TestSetup_Tracer_Close(t *testing.T) {
	assert := assert.New(t)
	err := Initialize("TEST", "0.0.0.0:6831")
	assert.NoError(err)
	tr := Tracer()
	assert.NotNil(tr)
	err = Close()
	assert.NoError(err)
}

func TestStartFinishHTTPSpan(t *testing.T) {
	assert := assert.New(t)
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	sp := StartHTTPSpan("/", req)
	assert.NotNil(sp)
	assert.IsType(&jaeger.Span{}, sp)
	jsp := sp.(*jaeger.Span)
	assert.NotNil(jsp)
	assert.Equal("HTTP GET /", jsp.OperationName())
	FinishHTTPSpan(jsp, 200)
	assert.NotNil(jsp)
}
