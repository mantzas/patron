package httprouter

import (
	"testing"

	"github.com/mantzas/patron/http"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	log.Setup(zerolog.DefaultFactory(log.DebugLevel))
	h := CreateHandler([]http.Route{http.NewRoute("/", "GET", nil)})
	assert.NotNil(h)
}
