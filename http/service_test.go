package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{}, false},
		{"failed with wrong option", args{[]Option{Ports(-1, -1)}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.options...)

			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestServer_ListenAndServer_Shutdown(t *testing.T) {
	assert := assert.New(t)
	s, err := New(Ports(10000, 10001))
	assert.NoError(err)
	go func() {
		s.ListenAndServe()
	}()
	err = s.shutdown()
	assert.NoError(err)
}
