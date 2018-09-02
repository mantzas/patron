package amqp

import (
	"testing"

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
	_, err = NewJSONMessage(make(chan bool))
	assert.Error(err)
}

func TestNewPublisher(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		url string
		exc string
		opt OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "fail, missing url", args: args{}, wantErr: true},
		{name: "fail, missing exchange", args: args{url: "url"}, wantErr: true},
		{name: "fail, missing exchange", args: args{url: "url", exc: "exc", opt: Timeout(0)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPublisher(tt.args.url, tt.args.exc, tt.args.opt)
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
