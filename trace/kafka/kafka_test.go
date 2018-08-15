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
		{name: "failure due to invalid data", data: make(chan bool), wantErr: true},
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
	seed := createKafkaBroker(t)
	got, err := NewAsyncProducer([]string{seed.Addr()})
	assert.NoError(err)
	assert.NotNil(got)
}

func TestAsyncProducer_SendMessage_Close(t *testing.T) {
	assert := assert.New(t)
	msg, err := NewJSONMessage("TOPIC", "TEST")
	assert.NoError(err)
	seed := createKafkaBroker(t)
	ap, err := NewAsyncProducer([]string{seed.Addr()})
	assert.NoError(err)
	assert.NotNil(ap)
	err = trace.Setup("test", "1.0.0", "0.0.0.0:6831", jaeger.SamplerTypeProbabilistic, 0.1)
	assert.NoError(err)
	_, ctx := trace.ChildSpan(context.Background(), "123", "cmp")
	err = ap.Send(ctx, msg)
	assert.NoError(err)
}

func createKafkaBroker(t *testing.T) *sarama.MockBroker {
	lead := sarama.NewMockBroker(t, 2)
	metadataResponse := new(sarama.MetadataResponse)
	metadataResponse.AddBroker(lead.Addr(), lead.BrokerID())
	metadataResponse.AddTopic("TOPIC", sarama.ErrNoError)
	metadataResponse.AddTopicPartition("TOPIC", 0, lead.BrokerID(), nil, nil, sarama.ErrNoError)

	prodSuccess := new(sarama.ProduceResponse)
	prodSuccess.Version = 11
	prodSuccess.AddTopicPartition("TOPIC", 0, sarama.ErrNoError)
	lead.Returns(prodSuccess)

	seed := sarama.NewMockBroker(t, 1)
	seed.Returns(metadataResponse)
	return seed
}
