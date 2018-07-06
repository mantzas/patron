package amqp

import (
	"context"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	assert := assert.New(t)
	m := NewMessage("xxx", []byte("test"))
	assert.Equal("xxx", m.contentType)
	assert.Equal([]byte("test"), m.body)
}

func TestNewJSONMessage(t *testing.T) {
	assert := assert.New(t)
	m, err := NewJSONMessage("xxx")
	assert.NoError(err)
	assert.Equal("application/json", m.contentType)
	assert.Equal([]byte(`"xxx"`), m.body)
}

func TestNewPublisher(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		url string
		exc string
	}
	tests := []struct {
		name          string
		args          args
		dialError     bool
		channelError  bool
		exchangeError bool
		wantErr       bool
	}{
		{name: "failure due to invalid url", args: args{url: "", exc: ""}, wantErr: true},
		{name: "failure due to invalid exchange", args: args{url: "url", exc: ""}, wantErr: true},
		{name: "failure due to dial", args: args{url: "url", exc: "exchange"}, dialError: true, wantErr: true},
		{name: "failure due to open channel", args: args{url: "url", exc: "exchange"}, channelError: true, wantErr: true},
		{name: "failure due to declare exchange", args: args{url: "url", exc: "exchange"}, exchangeError: true, wantErr: true},
		{name: "success", args: args{url: "url", exc: "exchange"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			cnn := &amqp.Connection{}
			chn := &amqp.Channel{}

			monkey.Patch(amqp.Dial, func(url string) (*amqp.Connection, error) {
				if tt.dialError {
					return nil, errors.New("DIAL ERROR")
				}
				return cnn, nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(cnn), "Channel", func(*amqp.Connection) (*amqp.Channel, error) {
				if tt.channelError {
					return nil, errors.New("CHANNEL ERROR")
				}
				return chn, nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(chn), "ExchangeDeclare", func(c *amqp.Channel, name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
				if tt.exchangeError {
					return errors.New("DECLARE EXCHANGE")
				}
				return nil
			})

			got, err := NewPublisher(tt.args.url, tt.args.exc)
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

func TestTracedPublisher_Close(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name       string
		closeError bool
		wantErr    bool
	}{
		{name: "failure closing channel", closeError: true, wantErr: true},
		{name: "success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			cnn := &amqp.Connection{}
			chn := &amqp.Channel{}

			monkey.Patch(amqp.Dial, func(url string) (*amqp.Connection, error) {
				return cnn, nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(cnn), "Channel", func(*amqp.Connection) (*amqp.Channel, error) {
				return chn, nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(chn), "ExchangeDeclare", func(c *amqp.Channel, name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
				return nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(chn), "Close", func(c *amqp.Channel) error {
				if tt.closeError {
					return errors.New("CHANNEL ERROR")
				}
				return nil
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(cnn), "Close", func(c *amqp.Connection) error {
				if tt.closeError {
					return errors.New("CONNECTION ERROR")
				}
				return nil
			})

			p, err := NewPublisher("XXX", "YYY")
			assert.NoError(err)
			err = p.Close(context.TODO())
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestTracedPublisher_Publish(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		publishError bool
		wantErr      bool
	}{
		//{name: "failure publishing", publishError: true, wantErr: true},
		{name: "success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			chn := &amqp.Channel{}
			monkey.PatchInstanceMethod(reflect.TypeOf(chn), "Publish", func(*amqp.Channel, string, string, bool, bool, amqp.Publishing) error {
				if tt.publishError {
					return errors.New("PUBLISH ERROR")
				}
				return nil
			})

			msg, err := NewJSONMessage("test")
			assert.NoError(err)
			tc := TracedPublisher{ch: chn}
			err = tc.Publish(context.TODO(), msg)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
