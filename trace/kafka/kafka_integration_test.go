// +build integration

package kafka

import (
	"context"
	"testing"
	"time"

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
	ch := make(chan *Result, 10)
	p, err := NewAsyncProducer(brokers, ch)
	assert.NoError(t, err)
	defer p.Close()
	err = p.Send(context.Background(), topic, payload)
	assert.NoError(t, err)
	res := <-ch
	assert.NoError(t, res.Err)
	assert.Equal(t, topic, res.Topic)
	assert.Equal(t, int32(0), res.Partition)
	assert.True(t, res.Offset > int64(0))
	err = p.SendRaw(context.Background(), topic, []byte(payload))
	assert.NoError(t, err)
	res = <-ch
	assert.NoError(t, res.Err)
	assert.Equal(t, topic, res.Topic)
	assert.Equal(t, int32(0), res.Partition)
	assert.True(t, res.Offset > int64(0))
}

var err error

func BenchmarkProducer_Send(b *testing.B) {
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	topic := "test-topic"
	payload := "TEST"
	brokers := []string{"localhost:9092"}
	p, err := NewProducer(brokers)
	assert.NoError(b, err)
	defer p.Close()
	err = p.Send(context.Background(), topic, payload)
	assert.NoError(b, err)
	time.Sleep(10 * time.Second)
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = p.Send(context.Background(), topic, payload)
	}
}

var res *Result

// TODO: this one never ends
func BenchmarkProducer_AsyncSend(b *testing.B) {
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	topic := "test-topic"
	payload := "TEST"
	brokers := []string{"localhost:9092"}
	ch := make(chan *Result, 10000000)
	go func() {
		for r := range ch {
			res = r
		}
	}()
	p, err := NewAsyncProducer(brokers, ch)
	assert.NoError(b, err)
	defer p.Close()
	err = p.Send(context.Background(), topic, payload)
	assert.NoError(b, err)
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = p.Send(context.Background(), topic, payload)
	}
}