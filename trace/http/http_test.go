package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mantzas/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestTracedClient_Do(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal("true", r.Header.Get("Mockpfx-Ids-Sampled"))
		assert.Equal("46", r.Header.Get("Mockpfx-Ids-Spanid"))
		assert.Equal("43", r.Header.Get("Mockpfx-Ids-Traceid"))
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	c := NewClient(1 * time.Second)
	req, err := http.NewRequest("GET", ts.URL, nil)
	assert.NoError(err)
	rsp, err := c.Do(context.Background(), req)
	assert.NoError(err)
	assert.NotNil(rsp)
	sp := mtr.FinishedSpans()[0]
	assert.NotNil(sp)
	assert.Equal(trace.HTTPOpName("GET", ts.URL), sp.OperationName)
}
