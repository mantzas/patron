package async

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFailureStrategy(t *testing.T) {
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
			c, err := New("test", proc.Process, &mockConsumerFactory{})
			assert.NoError(t, err)
			err = FailureStrategy(tt.args.fs)(c)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConsumerRetry(t *testing.T) {
	proc := mockProcessor{}
	type args struct {
		retries   int
		retryWait time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{retries: 3, retryWait: time.Second}, wantErr: false},
		{name: "invalid retries", args: args{retries: -1, retryWait: time.Second}, wantErr: true},
		{name: "invalid retry wait", args: args{retries: 3, retryWait: -1}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New("test", proc.Process, &mockConsumerFactory{})
			assert.NoError(t, err)
			err = ConsumerRetry(tt.args.retries, tt.args.retryWait)(c)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
