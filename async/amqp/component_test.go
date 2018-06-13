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
	type args struct {
		name  string
		url   string
		queue string
		p     async.Processor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"test", "url", "queue", &async.MockProcessor{}}, false},
		{"failed with invalid name", args{"", "url", "queue", &async.MockProcessor{}}, true},
		{"failed with invalid url", args{"test", "", "queue", &async.MockProcessor{}}, true},
		{"failed with invalid queue name", args{"test", "url", "", &async.MockProcessor{}}, true},
		{"failed with invalid processor", args{"test", "url", "queue", nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.name, tt.args.url, tt.args.queue, tt.args.p)
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
