package kafka

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/trace"
	"github.com/stretchr/testify/assert"
	jaeger "github.com/uber/jaeger-client-go"
)

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)
	m := NewMessage("TOPIC", []byte("TEST"))
	assert.Equal("TOPIC", m.topic)
	assert.Equal([]byte("TEST"), m.body)
}

func TestNewJSONMessage(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{name: "failure due to invalid data", data: make(chan bool, 0), wantErr: true},
		{name: "success", data: "TEST"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONMessage("TOPIC", tt.data)
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

func TestNewSyncProducer_Failure(t *testing.T) {
	assert := assert.New(t)
	got, err := NewAsyncProducer([]string{})
	assert.Error(err)
	assert.Nil(got)
}

func TestNewSyncProducer_Success(t *testing.T) {
	assert := assert.New(t)
	seed := createSeedBroker(t, false)
	got, err := NewAsyncProducer([]string{seed.Addr()})
	assert.NoError(err)
	assert.NotNil(got)
}

func TestAsyncProducer_SendMessage_Close(t *testing.T) {
	assert := assert.New(t)
	msg, err := NewJSONMessage("TOPIC", "TEST")
	assert.NoError(err)
	seed := createSeedBroker(t, true)
	ap, err := NewAsyncProducer([]string{seed.Addr()})
	assert.NoError(err)
	assert.NotNil(ap)
	trace.Setup("test", "0.0.0.0:6831", jaeger.SamplerTypeProbabilistic, 0.1)
	_, ctx := trace.StartChildSpan(context.Background(), "ttt", "cmp")
	ap.SendMessage(ctx, msg)
	assert.Error(<-ap.Error())
	ap.Close()
}

type mockAsyncProducer struct {
	in  chan *sarama.ProducerMessage
	err chan *sarama.ProducerError
}

func newMockAsyncProducer() *mockAsyncProducer {
	return &mockAsyncProducer{in: make(chan *sarama.ProducerMessage, 1), err: make(chan *sarama.ProducerError, 1)}
}

func (ma *mockAsyncProducer) AsyncClose() {
}

func (ma *mockAsyncProducer) Close() error {
	return nil
}

func (ma *mockAsyncProducer) Input() chan<- *sarama.ProducerMessage {
	return ma.in
}

func (ma *mockAsyncProducer) Successes() <-chan *sarama.ProducerMessage {
	return ma.in
}

func (ma *mockAsyncProducer) Errors() <-chan *sarama.ProducerError {
	return ma.err
}

func createSeedBroker(t *testing.T, retError bool) *sarama.MockBroker {
	seed := sarama.NewMockBroker(t, 1)
	lead := sarama.NewMockBroker(t, 2)

	metadataResponse := new(sarama.MetadataResponse)
	metadataResponse.AddBroker(lead.Addr(), lead.BrokerID())
	metadataResponse.AddTopicPartition("TOPIC", 0, lead.BrokerID(), nil, nil, sarama.ErrNoError)
	seed.Returns(metadataResponse)

	prodSuccess := new(sarama.ProduceResponse)
	if retError {
		prodSuccess.AddTopicPartition("TOPIC", 0, sarama.ErrDuplicateSequenceNumber)
	} else {
		prodSuccess.AddTopicPartition("TOPIC", 0, sarama.ErrNoError)
	}
	lead.Returns(prodSuccess)

	config := sarama.NewConfig()
	config.Producer.Flush.Messages = 10
	config.Producer.Return.Successes = true
	return seed
}
