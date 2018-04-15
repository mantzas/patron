package httprouter

import (
	"testing"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Setup(zerolog.DefaultFactory(log.DebugLevel))
}

func TestCreateHandler(t *testing.T) {
	assert := assert.New(t)
	h := CreateHandler([]patron.Route{patron.NewRoute("/", "GET", nil)})
	assert.NotNil(h)
}
