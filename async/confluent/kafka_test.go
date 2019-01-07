package confluent

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		name    string
		topics  []string
		brokers []string
		oo      []OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "failed, missing name", args: args{name: "", topics: []string{"topic"}, brokers: []string{"broker"}}, wantErr: true},
		{name: "failed, missing topic", args: args{name: "name", topics: []string{}, brokers: []string{"broker"}}, wantErr: true},
		{name: "failed, nil topic", args: args{name: "name", topics: nil, brokers: []string{"broker"}}, wantErr: true},
		{name: "failed, missing broker", args: args{name: "name", topics: []string{"topic"}, brokers: []string{}}, wantErr: true},
		{name: "failed, nil broker", args: args{name: "name", topics: []string{"topic"}, brokers: nil}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, "ct", tt.args.topics, tt.args.brokers, tt.args.oo...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func ErrorOption() OptionFunc {
	return func(c *consumer) error {
		return errors.New("TEST")
	}
}

func TestFactory_Create(t *testing.T) {

	expInfo := map[string]interface{}{
		"auto.offset.reset":               OffsetLatest,
		"bootstrap.servers":               "broker",
		"brokers":                         "broker",
		"buffer":                          1000,
		"client.id":                       "prometheus-name",
		"content-type":                    "ct",
		"go.application.rebalance.enable": true,
		"go.events.channel.enable":        true,
		"go.events.channel.size":          1000,
		"session.timeout.ms":              10000,
		"topics":                          "topic",
		"type":                            "kafka-consumer",
	}

	type fields struct {
		oo []OptionFunc
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "failed, error option", fields: fields{oo: []OptionFunc{ErrorOption()}}, wantErr: true},
		{name: "success", fields: fields{oo: []OptionFunc{}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := New("name", "ct", []string{"topic"}, []string{"broker"}, tt.fields.oo...)
			assert.NoError(t, err)
			got, err := f.Create()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, expInfo, got.Info())
			}
		})
	}
}
