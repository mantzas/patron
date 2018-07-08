package kafka

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/mantzas/patron/async"
	"github.com/mantzas/patron/encoding"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	proc := async.MockProcessor{}
	brokers := []string{"192.168.1.1"}
	topics := []string{"topic1"}
	type args struct {
		name     string
		proc     async.ProcessorFunc
		clientID string
		brokers  []string
		topics   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{name: "test", proc: proc.Process, clientID: "clID", brokers: brokers, topics: topics},
			wantErr: false,
		},
		{
			name:    "fails with missing name",
			args:    args{name: "", proc: proc.Process, clientID: "clID", brokers: brokers, topics: topics},
			wantErr: true,
		},
		{
			name:    "fails with missing processor",
			args:    args{name: "test", proc: nil, clientID: "clID", brokers: brokers, topics: topics},
			wantErr: true,
		},
		{
			name:    "fails with missing client id",
			args:    args{name: "test", proc: proc.Process, clientID: "", brokers: brokers, topics: topics},
			wantErr: true,
		},
		{
			name:    "fails with missing brokers",
			args:    args{name: "test", proc: proc.Process, clientID: "clID", brokers: []string{}, topics: topics},
			wantErr: true,
		},
		{
			name:    "fails with missing topics",
			args:    args{name: "test", proc: proc.Process, clientID: "clID", brokers: brokers, topics: []string{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.proc, tt.args.clientID, "", tt.args.brokers, tt.args.topics)
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
