package httprouter

import (
	"testing"

	patron_http "github.com/mantzas/patron/sync/http"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	h := CreateHandler([]patron_http.Route{patron_http.NewRoute("/", "GET", nil, false)})
	assert.NotNil(h)
}
