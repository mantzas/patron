// +build integration

package kafka

import (
	"context"
	"sync"
	"testing"

	"github.com/mantzas/patron/trace/kafka"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestConsume(t *testing.T) {
	mtr := mocktracer.New()
	opentracing.SetGlobalTracer(mtr)
	brokers := []string{"localhost:9092"}
	topics := []string{"test-topic"}
	// setup consumer
	f, err := New("test", topics, brokers)
	assert.NoError(t, err)
	cns, err := f.Create()
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, cns.Close())
	}()
	chMsg, chErr, err := cns.Consume(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, chMsg)
	assert.NotNil(t, chErr)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		//check send message
		m := <-chMsg
		assert.NotNil(t, m)
		assert.NotNil(t, m.Context())
		assert.NoError(t, m.Ack())
		assert.NoError(t, m.Nack())
	}()
	// setup producer
	p, err := kafka.NewProducer(brokers)
	assert.NoError(t, err)
	defer p.Close()
	err = p.Send(context.Background(), topics[0], "TEST")
	assert.NoError(t, err)
	wg.Wait()
}
