package amqp

import (
	"testing"

	"github.com/mantzas/patron/async"
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
		{"success", args{"url", "queue", &async.MockMesssageProcessor{}}, false},
		{"failed with invalid url", args{"", "queue", &async.MockMesssageProcessor{}}, true},
		{"failed with invalid queue name", args{"url", "", &async.MockMesssageProcessor{}}, true},
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
