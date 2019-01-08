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
	topic := "test-topic"
	payload := "TEST"
	brokers := []string{"localhost:9092"}
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	p, err := NewProducer(brokers)
	assert.NoError(t, err)
	defer p.Close()
	err = p.Send(context.Background(), topic, payload)
	assert.NoError(t, err)
	err = p.SendRaw(context.Background(), topic, []byte(payload))
	assert.NoError(t, err)
}

func TestAsyncSend(t *testing.T) {
	topic := "test-topic"
	payload := "TEST"
	brokers := []string{"localhost:9092"}
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	p, err := NewAsyncProducer(brokers)
	assert.NoError(t, err)
	defer p.Close()
	err = p.Send(context.Background(), topic, payload)
	assert.NoError(t, err)
	res := <-p.Results()
	assert.NoError(t, res.Err)
	assert.Equal(t, topic, res.Topic)
	assert.Equal(t, int32(0), res.Partition)
	assert.True(t, res.Offset > int64(0))
	err = p.SendRaw(context.Background(), topic, []byte(payload))
	assert.NoError(t, err)
	res = <-p.Results()
	assert.NoError(t, res.Err)
	assert.Equal(t, topic, res.Topic)
	assert.Equal(t, int32(0), res.Partition)
	assert.True(t, res.Offset > int64(0))
}
