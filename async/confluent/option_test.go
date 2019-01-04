package confluent

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
		errMsg  string
	}{
		{name: "success", args: args{cfg: map[string]interface{}{"test": "test"}}, wantErr: false},
		{name: "failure, nil config", args: args{cfg: nil}, wantErr: true, errMsg: "config is nil"},
		{name: "failure, empty config", args: args{cfg: make(map[string]interface{})}, wantErr: true, errMsg: "config is empty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &kafka.ConfigMap{
				"auto.offset.reset": OffsetLatest,
			}
			c := consumer{cfg: cfg}
			err := Config(tt.args.cfg)(&c)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
