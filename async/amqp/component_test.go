package amqp

import (
	"context"
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
		url   string
		queue string
		p     async.Processor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"url", "queue", &async.MockProcessor{}}, false},
		{"failed with invalid url", args{"", "queue", &async.MockProcessor{}}, true},
		{"failed with invalid queue name", args{"url", "", &async.MockProcessor{}}, true},
		{"failed with invalid processor", args{"url", "queue", nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.url, tt.args.queue, tt.args.p)
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

func Test_createContext(t *testing.T) {
	assert := assert.New(t)
	tbl := amqp.Table{}
	tbl["key1"] = "val1"
	tbl["key2"] = "val2"
	ctx, cnl := createContext(context.Background(), tbl)
	assert.NotNil(cnl)
	assert.Equal(tbl["key1"], ctx.Value(amqpContextKey("key1")))
	assert.Equal(tbl["key2"], ctx.Value(amqpContextKey("key2")))
}
