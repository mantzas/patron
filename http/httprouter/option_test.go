package httprouter

import (
	"testing"

	"github.com/mantzas/patron/http"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	assert := assert.New(t)
	s, err := http.New("test", []http.Route{http.NewRoute("/", "GET", nil)}, Handler())
	assert.NoError(err)
	assert.NotNil(s)
}
