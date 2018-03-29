package httprouter

import (
	"testing"

	"github.com/mantzas/patron/http"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	h := CreateHandler([]http.Route{http.NewRoute("/", "GET", nil)})
	assert.NotNil(h)
}
