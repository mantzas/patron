package patron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTracing(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		addr         string
		name         string
		samplerType  string
		samplerParam float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failure due to missing sampler type", args{addr: "0.0.0.0:6831", name: "TEST", samplerType: "", samplerParam: 1}, true},
		{"failure due to empty agent address", args{addr: "", name: "TEST", samplerType: "const", samplerParam: 1}, true},
		{"failure due to service name missing", args{addr: "0.0.0.0:6831", name: "", samplerType: "const", samplerParam: 1}, true},
		{"success", args{addr: "0.0.0.0:6831", name: "TEST", samplerType: "const", samplerParam: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{name: tt.args.name}
			err := Tracing(tt.args.addr, tt.args.samplerType, tt.args.samplerParam)(&s)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
