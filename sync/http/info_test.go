package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/info"
	"github.com/stretchr/testify/assert"
)

func Test_Route(t *testing.T) {
	info.AddName("Name")
	info.AddMetric("Name", "Description")
	r := infoRoute()
	mux := http.NewServeMux()
	mux.HandleFunc(r.Pattern, r.Handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/info")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, json.TypeCharset, resp.Header.Get(encoding.ContentTypeHeader))
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"name":"Name","metrics":[{"name":"Name","description":"Description"}]}`, string(body))
}
