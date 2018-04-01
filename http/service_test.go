package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testRoutes = []Route{NewRoute("/", "Get", func(w http.ResponseWriter, r *http.Request) {})}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		name    string
		routes  []Route
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success", args{"test", testRoutes, []Option{}}, false},
		{"failed with missing name", args{"", testRoutes, []Option{}}, true},
		{"failed with missing routes", args{"test", nil, []Option{}}, true},
		{"failed with wrong option", args{"test", testRoutes, []Option{Ports(-1, -1)}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.name, tt.args.routes, tt.args.options...)
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
	s, err := New("test", testRoutes, Ports(10000, 10001))
	assert.NoError(err)
	go func() {
		s.Run()
	}()
	err = s.shutdown()
	assert.NoError(err)
}
