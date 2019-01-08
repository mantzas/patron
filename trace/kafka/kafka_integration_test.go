// +build integration

package kafka

import (
	"context"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	p, err := NewProducer([]string{"localhost:9092"})
	assert.NoError(t, err)
	defer p.Close()
	msg := NewMessage("test-topic", []byte("TEST"))
	err = p.Send(context.Background(), msg)
	assert.NoError(t, err)
}

func TestAsyncSend(t *testing.T) {
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	p, err := NewAsyncProducer([]string{"localhost:9092"})
	assert.NoError(t, err)
	defer p.Close()
	msg := NewMessage("test-topic", []byte("TEST"))
	err = p.Send(context.Background(), msg)
	assert.NoError(t, err)
	res := <-p.Results()
	assert.NoError(t, res.Err)
	assert.Equal(t, "test-topic", res.Topic)
	assert.Equal(t, int32(0), res.Partition)
	assert.True(t, res.Offset > int64(0))
}
