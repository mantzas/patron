package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		name    string
		h       http.Handler
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"test", http.DefaultServeMux, []Option{}}, false},
		{"failed with missing name", args{"", http.DefaultServeMux, []Option{}}, true},
		{"failed with missing handler", args{"test", nil, []Option{}}, true},
		{"failed with wrong option", args{"test", http.DefaultServeMux, []Option{Ports(-1, -1)}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.name, tt.args.h, tt.args.options...)
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
	s, err := New("test", http.DefaultServeMux, Ports(10000, 10001))
	assert.NoError(err)
	go func() {
		s.ListenAndServe()
	}()
	err = s.shutdown()
	assert.NoError(err)
}
