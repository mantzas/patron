package amqp

import (
	"context"
	"reflect"
	"testing"

	"github.com/bouk/monkey"
	"github.com/mantzas/patron/trace"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	jaeger "github.com/uber/jaeger-client-go"
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
	_, err = NewJSONMessage(make(chan bool))
	assert.Error(err)
}

func TestNewPublisher(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		url string
		exc string
	}
	validArgs := args{url: "url", exc: "exchange"}
	tests := []struct {
		name          string
		args          args
		dialError     bool
		channelError  bool
		exchangeError bool
		wantErr       bool
	}{
		{name: "failure due to invalid url", args: args{}, wantErr: true},
		{name: "failure due to invalid exchange", args: args{url: "url"}, wantErr: true},
		{name: "failure due to dial", args: validArgs, dialError: true, wantErr: true},
		{name: "failure due to open channel", args: validArgs, channelError: true, wantErr: true},
		{name: "failure due to declare exchange", args: validArgs, exchangeError: true, wantErr: true},
		{name: "success", args: validArgs},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			cnn := &amqp.Connection{}
			chn := &amqp.Channel{}

			monkey.Patch(amqp.Dial, func(string) (*amqp.Connection, error) {
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

			monkey.PatchInstanceMethod(reflect.TypeOf(chn), "ExchangeDeclare", func(
				*amqp.Channel,
				string,
				string,
				bool,
				bool,
				bool,
				bool,
				amqp.Table,
			) error {
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
			var nilErr error

			monkey.Patch(amqp.Dial, func(string) (*amqp.Connection, error) {
				return cnn, nilErr
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(cnn), "Channel", func(*amqp.Connection) (*amqp.Channel, error) {
				return chn, nilErr
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(chn), "ExchangeDeclare", func(
				*amqp.Channel,
				string,
				string,
				bool,
				bool,
				bool,
				bool,
				amqp.Table,
			) error {
				return nilErr
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(chn), "Close", func(*amqp.Channel) error {
				if tt.closeError {
					return errors.New("CHANNEL ERROR")
				}
				return nilErr
			})

			monkey.PatchInstanceMethod(reflect.TypeOf(cnn), "Close", func(*amqp.Connection) error {
				if tt.closeError {
					return errors.New("CONNECTION ERROR")
				}
				return nilErr
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
	err := trace.Setup("test", "0.0.0.0:6831", jaeger.SamplerTypeProbabilistic, 0.1)
	assert.NoError(err)
	_, ctx := trace.StartChildSpan(context.Background(), "ttt", "cmp")
	tests := []struct {
		name         string
		publishError bool
		wantErr      bool
	}{
		{name: "failure publishing", publishError: true, wantErr: true},
		{name: "success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer monkey.UnpatchAll()
			chn := &amqp.Channel{}
			monkey.PatchInstanceMethod(reflect.TypeOf(chn), "Publish", func(
				*amqp.Channel,
				string,
				string,
				bool,
				bool,
				amqp.Publishing,
			) error {
				if tt.publishError {
					return errors.New("PUBLISH ERROR")
				}
				return nil
			})

			msg, err := NewJSONMessage("test")
			assert.NoError(err)
			tc := TracedPublisher{ch: chn}
			err = tc.Publish(ctx, msg)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
