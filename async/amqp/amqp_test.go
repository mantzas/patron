package amqp

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	type args struct {
		url      string
		queue    string
		exchange string
		opt      OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{url: "amqp://guest:guest@localhost:5672/", queue: "queue", exchange: "exchange", opt: Buffer(100)}, false},
		{"fail, invalid url", args{url: "", queue: "queue", exchange: "exchange", opt: Buffer(100)}, true},
		{"fail, invalid queue name", args{url: "url", queue: "", exchange: "exchange", opt: Buffer(100)}, true},
		{"fail, invalid exchange name", args{url: "url", queue: "queue", exchange: "", opt: Buffer(100)}, true},
		{"fail, invalid opt", args{url: "url", queue: "queue", exchange: "exchange", opt: Buffer(-100)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConsumer(tt.args.url, tt.args.queue, tt.args.exchange, tt.args.opt)
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

func Test_mapHeader(t *testing.T) {
	hh := amqp.Table{"test1": 10, "test2": 0.11}
	mm := map[string]string{"test1": "10", "test2": "0.11"}
	assert.Equal(t, mm, mapHeader(hh))
}

func TestConsumer_Info(t *testing.T) {
	c, err := NewConsumer("url", "queue", "exchange")
	assert.NoError(t, err)
	expected := make(map[string]interface{})
	expected["type"] = "amqp-consumer"
	expected["queue"] = "queue"
	expected["exchange"] = "exchange"
	expected["requeue"] = true
	expected["buffer"] = 1000
	expected["url"] = "url"
	assert.Equal(t, expected, c.Info())
}
