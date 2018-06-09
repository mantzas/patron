package patron

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go"
)

func TestTracing(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		sampler  jaeger.Sampler
		reporter jaeger.Reporter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failure due to nil sampler", args{nil, jaeger.NewNullReporter()}, true},
		{"failure due to nil reporter", args{jaeger.NewConstSampler(true), nil}, true},
		{"success", args{jaeger.NewConstSampler(true), jaeger.NewNullReporter()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{}
			err := Tracing(tt.args.sampler, tt.args.reporter)(&s)

			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
