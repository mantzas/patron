package confluent

import (
	"testing"

	"github.com/mantzas/patron/errors"
	"github.com/stretchr/testify/assert"
)

func ErrorOption() OptionFunc {
	return func(k *KafkaProducer) error {
		return errors.New("TEST")
	}
}

func TestNewMessage(t *testing.T) {
	topic := "xxx"
	body := []byte("Test")
	msg := NewMessage(topic, body)
	assert.Equal(t, topic, msg.topic)
	assert.Equal(t, body, msg.body)
}

func TestNewJSONMessage(t *testing.T) {
	topic := "xxx"
	body := struct {
		ID   int
		Name string
	}{
		ID:   1,
		Name: "Test",
	}
	expected := []byte(`{"ID":1,"Name":"Test"}`)
	msg, err := NewJSONMessage(topic, body)
	assert.NoError(t, err)
	assert.Equal(t, topic, msg.topic)
	assert.Equal(t, expected, msg.body)
}

func TestNewJSONMessageError(t *testing.T) {
	topic := "xxx"
	body := make(chan bool)
	msg, err := NewJSONMessage(topic, body)
	assert.Error(t, err)
	assert.Nil(t, msg)
}

func TestNewProducer(t *testing.T) {
	type args struct {
		brokers []string
		oo      []OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "failed, no brokers", args: args{}, wantErr: true},
		{name: "failed, invalid option", args: args{brokers: []string{"xxx"}, oo: []OptionFunc{ErrorOption()}}, wantErr: true},
		{name: "success", args: args{brokers: []string{"xxx"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProducer(tt.args.brokers, tt.args.oo...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Nil(t, got.Error())
			}
		})
	}
}

func TestNewAsyncProducer(t *testing.T) {
	type args struct {
		brokers []string
		oo      []OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "failed, no brokers", args: args{}, wantErr: true},
		{name: "failed, invalid option", args: args{brokers: []string{"xxx"}, oo: []OptionFunc{ErrorOption()}}, wantErr: true},
		{name: "success", args: args{brokers: []string{"xxx"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAsyncProducer(tt.args.brokers, tt.args.oo...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.NotNil(t, got.Error())
			}
		})
	}
}
