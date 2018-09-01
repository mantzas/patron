package async

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailureStrategy(t *testing.T) {
	assert := assert.New(t)
	proc := mockProcessor{}
	type args struct {
		fs FailStrategy
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{fs: NackExitStrategy}, wantErr: false},
		{name: "invalid strategy (lower)", args: args{fs: -1}, wantErr: true},
		{name: "invalid strategy (higher)", args: args{fs: 3}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(proc.Process, &mockConsumer{})
			assert.NoError(err)
			err = FailureStrategy(tt.args.fs)(c)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}
