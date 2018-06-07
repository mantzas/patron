package http

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func ErrorOption() Option {
	return func(s *Component) error {
		return errors.New("TEST")
	}
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	reporter := jaeger.NewInMemoryReporter()
	tr, trCloser := jaeger.NewTracer("test", jaeger.NewConstSampler(true), reporter)
	defer trCloser.Close()
	tests := []struct {
		name    string
		tr      opentracing.Tracer
		hg      handlerGen
		options []Option
		wantErr bool
	}{
		{"success with no options", tr, testCreateHandler, []Option{}, false},
		{"success with options", tr, testCreateHandler, []Option{Port(50000)}, false},
		{"failed with error option", tr, testCreateHandler, []Option{ErrorOption()}, true},
		{"failed with missing tracer", nil, nil, []Option{}, true},
		{"failed with missing handler gen", nil, testCreateHandler, []Option{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.tr, tt.hg, tt.options...)
			if tt.wantErr {
				assert.Error(err)
				assert.Nil(got)
			} else {
				assert.NoError(err)
				assert.NotNil(got)
			}
		})
	}
}

func TestComponent_ListenAndServer_DefaultRoutes_Shutdown(t *testing.T) {
	assert := assert.New(t)
	reporter := jaeger.NewInMemoryReporter()
	tr, trCloser := jaeger.NewTracer("test", jaeger.NewConstSampler(true), reporter)
	defer trCloser.Close()
	s, err := New(tr, testCreateHandler)
	assert.NoError(err)
	go func() {
		err = s.Run(context.TODO())
		assert.NoError(err)
	}()
	assert.Len(s.routes, 11)
	err = s.Shutdown(context.TODO())
	assert.NoError(err)
}

func testCreateHandler(routes []Route) http.Handler {
	return http.NewServeMux()
}

func Test_createHTTPServer(t *testing.T) {
	assert := assert.New(t)
	s := createHTTPServer(10000, nil)
	assert.NotNil(s)
	assert.Equal(":10000", s.Addr)
	assert.Equal(5*time.Second, s.ReadTimeout)
	assert.Equal(60*time.Second, s.WriteTimeout)
	assert.Equal(120*time.Second, s.IdleTimeout)
}
