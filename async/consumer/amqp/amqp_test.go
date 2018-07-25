package amqp

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		name   string
		url    string
		queue  string
		buffer int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{name: "test", url: "url", queue: "queue", buffer: 10}, false},
		{"failed with invalid name", args{name: "", url: "url", queue: "queue", buffer: 10}, true},
		{"failed with invalid url", args{name: "test", url: "", queue: "queue", buffer: 10}, true},
		{"failed with invalid queue name", args{name: "test", url: "url", queue: "", buffer: 10}, true},
		{"failed with invalid buffer", args{name: "test", url: "url", queue: "queue", buffer: -10}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.url, tt.args.queue, true, tt.args.buffer)
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

func Test_mapHeader(t *testing.T) {
	assert := assert.New(t)
	hh := amqp.Table{"test1": 10, "test2": 0.11}
	mm := map[string]string{"test1": "10", "test2": "0.11"}
	assert.Equal(mm, mapHeader(hh))
}
