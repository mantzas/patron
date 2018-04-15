package amqp

import (
	"context"
	"testing"

	"github.com/mantzas/patron/async"
	"github.com/stretchr/testify/assert"
)

type testMesssageProcessor struct {
}

func (tmp testMesssageProcessor) Process(ctx context.Context, msg []byte) error {

	return nil
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		url   string
		queue string
		mp    async.MessageProcessor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"url", "queue", &testMesssageProcessor{}}, false},
		{"failed with invalid url", args{"", "queue", &testMesssageProcessor{}}, true},
		{"failed with invalid queue name", args{"url", "", &testMesssageProcessor{}}, true},
		{"failed with invalid processor", args{"url", "queue", nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.url, tt.args.queue, tt.args.mp)
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
