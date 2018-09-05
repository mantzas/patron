package kafka

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		buf int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{buf: 100}, wantErr: false},
		{name: "invalid buffer", args: args{buf: -100}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Consumer{}
			err := Buffer(tt.args.buf)(&c)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestStart(t *testing.T) {
	assert := assert.New(t)
	c := Consumer{}
	err := Start(1000)(&c)
	assert.NoError(err)
}

func TestTimeout(t *testing.T) {
	assert := assert.New(t)
	c := Consumer{cfg: sarama.NewConfig()}
	err := Timeout(time.Second)(&c)
	assert.NoError(err)
}

func TestVersion(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{version: "1.0.0"}, wantErr: false},
		{name: "failed due to empty", args: args{version: ""}, wantErr: true},
		{name: "failed due to invalid", args: args{version: "1.0.0.0"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New("test", "", "topic", []string{"test"})
			assert.NoError(err)
			err = Version(tt.args.version)(c)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
