package http

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func ErrorOption() OptionFunc {
	return func(s *Component) error {
		return errors.New("TEST")
	}
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		options []OptionFunc
		wantErr bool
	}{
		{"success with no options", []OptionFunc{}, false},
		{"success with options", []OptionFunc{Port(50000)}, false},
		{"failed with error option", []OptionFunc{ErrorOption()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.options...)
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

func TestComponent_ListenAndServe_DefaultRoutes_Shutdown(t *testing.T) {
	assert := assert.New(t)
	rr := []Route{NewRoute("/", "GET", nil, true)}
	s, err := New(Routes(rr))
	assert.NoError(err)
	go func() {
		assert.Error(s.Run(context.Background()))
	}()
	time.Sleep(100 * time.Millisecond)
	assert.Len(s.routes, 13)
	assert.NoError(s.Shutdown(context.Background()))
}

func TestComponent_ListenAndServeTLS_DefaultRoutes_Shutdown(t *testing.T) {
	assert := assert.New(t)
	rr := []Route{NewRoute("/", "GET", nil, true)}
	s, err := New(Routes(rr), Secure("testdata/server.pem", "testdata/server.pem"))
	assert.NoError(err)
	go func() {
		assert.Error(s.Run(context.Background()))
	}()
	time.Sleep(100 * time.Millisecond)
	assert.Len(s.routes, 13)
	assert.NoError(s.Shutdown(context.Background()))
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

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	h := createHandler([]Route{NewRoute("/", "GET", nil, false)})
	assert.NotNil(h)
}
