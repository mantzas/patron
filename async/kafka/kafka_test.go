package kafka

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/encoding"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	brokers := []string{"192.168.1.1"}
	type args struct {
		name     string
		clientID string
		brokers  []string
		topic    string
		buffer   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{name: "test", clientID: "clID", brokers: brokers, topic: "topic1", buffer: 0},
			wantErr: false,
		},
		{
			name:    "fails with missing name",
			args:    args{name: "", clientID: "clID", brokers: brokers, topic: "topic1", buffer: 0},
			wantErr: true,
		},
		{
			name:    "fails with missing client id",
			args:    args{name: "test", clientID: "", brokers: brokers, topic: "topic1", buffer: 0},
			wantErr: true,
		},
		{
			name:    "fails with missing brokers",
			args:    args{name: "test", clientID: "clID", brokers: []string{}, topic: "topic1", buffer: 0},
			wantErr: true,
		},
		{
			name:    "fails with missing topics",
			args:    args{name: "test", clientID: "clID", brokers: brokers, topic: "", buffer: 0},
			wantErr: true,
		},
		{
			name:    "fails with invalid buffer",
			args:    args{name: "test", clientID: "clID", brokers: brokers, topic: "topic1", buffer: -1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.clientID, "", tt.args.topic, tt.args.brokers, tt.args.buffer, OffsetOldest)
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

func Test_determineContentType(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		hdr []*sarama.RecordHeader
	}

	validHdr := &sarama.RecordHeader{
		Key:   []byte(encoding.ContentTypeHeader),
		Value: []byte("val1"),
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"failure", args{hdr: []*sarama.RecordHeader{}}, "", true},
		{"success", args{hdr: []*sarama.RecordHeader{validHdr}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := determineContentType(tt.args.hdr)
			if tt.wantErr {
				assert.Empty(got)
				assert.Error(err)
			} else {
				assert.NotNil(got)
				assert.NoError(err)
			}
		})
	}
}

// func TestRun_Shutdown(t *testing.T) {
// 	assert := assert.New(t)
// 	br := createSeedBroker(t, false)
// 	c, err := New("test", "1", "12", []string{br.Addr()}, "TOPIC", 0)
// 	assert.NoError(err)
// 	assert.NotNil(c)
// 	go func() {
// 		c.Consume(context.Background())
// 	}()
// 	time.Sleep(100 * time.Millisecond)
// 	assert.NoError(c.Close())
// }

// func createSeedBroker(t *testing.T, retError bool) *sarama.MockBroker {
// 	seed := sarama.NewMockBroker(t, 1)
// 	lead := sarama.NewMockBroker(t, 2)

// 	metadataResponse := new(sarama.MetadataResponse)
// 	metadataResponse.AddBroker(lead.Addr(), lead.BrokerID())
// 	metadataResponse.AddTopicPartition("TOPIC", 0, lead.BrokerID(), nil, nil, sarama.ErrNoError)
// 	seed.Returns(metadataResponse)

// 	prodSuccess := new(sarama.ProduceResponse)
// 	if retError {
// 		prodSuccess.AddTopicPartition("TOPIC", 0, sarama.ErrDuplicateSequenceNumber)
// 	} else {
// 		prodSuccess.AddTopicPartition("TOPIC", 0, sarama.ErrNoError)
// 	}
// 	lead.Returns(prodSuccess)

// 	config := sarama.NewConfig()
// 	config.Producer.Flush.Messages = 10
// 	config.Producer.Return.Successes = true
// 	return seed
// }
