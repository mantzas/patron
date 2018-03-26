package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPorts(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		port      int
		pprofPort int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{30000, 30001}, false},
		{"error for same port", args{30000, 30000}, true},
		{"error for port number out of range", args{-1, 30001}, true},
		{"error for pprof port number out of range", args{30000, -1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, err := New(Ports(tt.args.port, tt.args.pprofPort))

			if tt.wantErr {
				assert.Nil(s)
				assert.Error(err)
			} else {
				assert.NotNil(s)
				assert.NoError(err)
			}
		})
	}
}
