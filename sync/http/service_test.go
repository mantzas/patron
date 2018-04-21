package http

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		hg      handlerGen
		options []Option
		wantErr bool
	}{
		{"success with no options", testCreateHandler, []Option{}, false},
		{"success with options", testCreateHandler, []Option{SetPorts(50000)}, false},
		{"failed with missing handler gen", nil, []Option{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.hg, tt.options...)
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

func TestService_ListenAndServer_Shutdown(t *testing.T) {
	assert := assert.New(t)
	s, err := New(testCreateHandler)
	assert.NoError(err)
	go func() {
		err = s.Run(context.TODO())
		assert.NoError(err)
	}()
	err = s.Shutdown(context.TODO())
	assert.NoError(err)
}

func testCreateHandler(routes []Route) http.Handler {
	return http.NewServeMux()
}
