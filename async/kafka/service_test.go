package kafka

import (
	"testing"

	"github.com/mantzas/patron/async"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	type args struct {
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
		{"success", args{&async.MockMesssageProcessor{}, "clID", []string{"192.168.1.1"}, []string{"topic1"}}, false},
		{"fails with missing processor", args{nil, "clID", []string{"192.168.1.1"}, []string{"topic1"}}, true},
		{"fails with missing client id", args{&async.MockMesssageProcessor{}, "", []string{"192.168.1.1"}, []string{"topic1"}}, true},
		{"fails with missing brokers", args{&async.MockMesssageProcessor{}, "clID", []string{}, []string{"topic1"}}, true},
		{"fails with missing topics", args{&async.MockMesssageProcessor{}, "clID", []string{"192.168.1.1"}, []string{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.p, tt.args.clientID, tt.args.brokers, tt.args.topics)
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
