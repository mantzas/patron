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
	type args struct {
		name     string
		p        async.Processor
		clientID string
		brokers  []string
		topics   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"test", &async.MockProcessor{}, "clID", []string{"192.168.1.1"}, []string{"topic1"}}, false},
		{"fails with missing name", args{"", &async.MockProcessor{}, "clID", []string{"192.168.1.1"}, []string{"topic1"}}, true},
		{"fails with missing processor", args{"test", nil, "clID", []string{"192.168.1.1"}, []string{"topic1"}}, true},
		{"fails with missing client id", args{"test", &async.MockProcessor{}, "", []string{"192.168.1.1"}, []string{"topic1"}}, true},
		{"fails with missing brokers", args{"test", &async.MockProcessor{}, "clID", []string{}, []string{"topic1"}}, true},
		{"fails with missing topics", args{"test", &async.MockProcessor{}, "clID", []string{"192.168.1.1"}, []string{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.p, tt.args.clientID, tt.args.brokers, tt.args.topics)
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
		{"failure", args{[]*sarama.RecordHeader{}}, "", true},
		{"success", args{[]*sarama.RecordHeader{validHdr}}, "", false},
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
