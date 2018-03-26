package pprof

import (
	"context"
	_ "net/http/pprof"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	s := New(10001)
	assert.NotNil(s.srv)
	assert.Equal(":10001", s.srv.Addr)
}

func TestServer_SetPort(t *testing.T) {
	assert := assert.New(t)
	s := New(10001)
	s.SetPort(10002)
	assert.NotNil(s.srv)
	assert.Equal(":10002", s.srv.Addr)
}

func TestServer_GetAddr(t *testing.T) {
	assert := assert.New(t)
	s := New(10001)
	s.SetPort(10002)
	assert.NotNil(s.srv)
	assert.Equal(":10002", s.GetAddr())
}

func TestServer_ListenAndServer_Shutdown(t *testing.T) {
	assert := assert.New(t)
	s := New(10001)
	go func() {
		s.ListenAndServe()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := s.Shutdown(ctx)
	assert.NoError(err)
}
