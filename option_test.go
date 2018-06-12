package patron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracing(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		addr string
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failure due to empty agent address", args{"", "TEST"}, true},
		{"failure due to service name missing", args{"0.0.0.0:6831", ""}, true},
		{"success", args{"0.0.0.0:6831", "TEST"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{name: tt.args.name}
			err := Tracing(tt.args.addr)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
