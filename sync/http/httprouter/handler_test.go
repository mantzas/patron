package httprouter

import (
	"testing"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/mantzas/patron/sync/http"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Setup(zerolog.DefaultFactory(log.DebugLevel))
}

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	h := CreateHandler([]http.Route{http.NewRoute("/", "GET", nil)})
	assert.NotNil(h)
}
