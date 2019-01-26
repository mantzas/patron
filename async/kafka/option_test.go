package kafka

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	type args struct {
		cfg map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantAck bool
		errMsg  string
	}{
		{
			name:    "success",
			args:    args{cfg: map[string]interface{}{"test": "test"}},
			wantAck: false,
			wantErr: false,
		},
		{
			name:    "success, manual message ack",
			args:    args{cfg: map[string]interface{}{"enable.auto.commit": false}},
			wantAck: true,
			wantErr: false,
		},
		{
			name:    "failure, nil config",
			args:    args{cfg: nil},
			wantAck: false,
			wantErr: true,
			errMsg:  "config is nil",
		},
		{
			name:    "failure, empty config",
			args:    args{cfg: make(map[string]interface{})},
			wantAck: false,
			wantErr: true,
			errMsg:  "config is empty",
		},
		{
			name:    "failure, enable.auto.commit not bool",
			args:    args{cfg: map[string]interface{}{"enable.auto.commit": 1}},
			wantAck: false,
			wantErr: true,
			errMsg:  "enable.auto.commit should be boolean",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &kafka.ConfigMap{
				"auto.offset.reset": "latest",
			}
			c := consumer{cfg: cfg}
			err := Config(tt.args.cfg)(&c)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAck, c.ack)
			}
		})
	}
}
