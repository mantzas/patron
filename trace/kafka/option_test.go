package kafka

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
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
			cfg := &kafka.ConfigMap{}
			c := Producer{cfg: cfg}
			err := Config(tt.args.cfg)(&c)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		enc encoding.EncodeFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "failed, nil encoder", args: args{enc: nil}, wantErr: true},
		{name: "success", args: args{enc: json.Encode}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &kafka.ConfigMap{}
			c := Producer{cfg: cfg}
			err := Encode(tt.args.enc)(&c)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
