package amqp

import (
	"context"
	"testing"

	"github.com/mantzas/patron/worker"
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
		mp    worker.MessageProcessor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"url", "queue", &testMesssageProcessor{}}, false},
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
