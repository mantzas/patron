package httprouter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mantzas/patron/sync"
	patron_http "github.com/mantzas/patron/sync/http"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	h := CreateHandler([]patron_http.Route{patron_http.NewRoute("/", "GET", nil, nil, false)})
	assert.NotNil(h)
}

func TestParamExtractor(t *testing.T) {
	assert := assert.New(t)
	req, err := http.NewRequest(http.MethodGet, "/users/1/status", nil)
	assert.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	var fields map[string]string

	proc := func(_ context.Context, req *sync.Request) (*sync.Response, error) {
		fields = req.Fields
		return nil, nil
	}

	h := CreateHandler([]patron_http.Route{patron_http.NewRoute("/users/:id/status", "GET", proc, ParamExtractor, false)})
	h.ServeHTTP(httptest.NewRecorder(), req)
	assert.Equal("1", fields["id"])
}
