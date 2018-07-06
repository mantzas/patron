package amqp

import (
	"testing"

	"github.com/mantzas/patron/async"
	agr_errors "github.com/mantzas/patron/errors"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	proc := async.MockProcessor{}
	type args struct {
		name  string
		url   string
		queue string
		proc  async.ProcessorFunc
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{name: "test", url: "url", queue: "queue", proc: proc.Process}, false},
		{"failed with invalid name", args{name: "", url: "url", queue: "queue", proc: proc.Process}, true},
		{"failed with invalid url", args{name: "test", url: "", queue: "queue", proc: proc.Process}, true},
		{"failed with invalid queue name", args{name: "test", url: "url", queue: "", proc: proc.Process}, true},
		{"failed with invalid processor", args{name: "test", url: "url", queue: "queue", proc: nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.url, tt.args.queue, tt.args.proc)
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

func Test_handlerMessageError(t *testing.T) {
	assert := assert.New(t)
	agr := agr_errors.New()
	d := &amqp.Delivery{
		MessageId: "1",
	}
	handlerMessageError(d, agr, errors.New("test"), "message")
	assert.Equal(2, agr.Count())
	assert.Equal("message: test\nfailed to NACK message 1: delivery not initialized\n", agr.Error())
}
