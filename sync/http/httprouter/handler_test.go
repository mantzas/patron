package httprouter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	patron_http "github.com/mantzas/patron/sync/http"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	h := CreateHandler([]patron_http.Route{patron_http.NewRoute("/", "GET", nil)})
	assert.NotNil(h)
}

func Test_PprofHandlers(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(CreateHandler([]patron_http.Route{}))
	defer server.Close()

	tests := []struct {
		name string
		path string
		want int
	}{
		{"index", "/debug/pprof/", 200},
		{"cmdline", "/debug/pprof/cmdline/", 200},
		{"profile", "/debug/pprof/profile/?seconds=1", 200},
		{"symbol", "/debug/pprof/symbol/", 200},
		{"trace", "/debug/pprof/trace/?seconds=0.1", 200},
		{"heap", "/debug/pprof/heap/", 200},
		{"goroutine", "/debug/pprof/goroutine/", 200},
		{"block", "/debug/pprof/block/", 200},
		{"threadcreate", "/debug/pprof/threadcreate/", 200},
		{"mutex", "/debug/pprof/mutex/", 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(fmt.Sprintf("%s/%s", server.URL, tt.path))
			assert.NoError(err)
			assert.Equal(tt.want, resp.StatusCode)
		})
	}
}
