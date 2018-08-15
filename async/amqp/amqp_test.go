package amqp

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		url      string
		queue    string
		exchange string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{url: "url", queue: "queue", exchange: "exchange"}, false},
		{"fail, invalid url", args{url: "", queue: "queue", exchange: "exchange"}, true},
		{"fail, invalid queue name", args{url: "url", queue: "", exchange: "exchange"}, true},
		{"fail, invalid queue name", args{url: "url", queue: "queue", exchange: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.url, tt.args.queue, tt.args.exchange)
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
