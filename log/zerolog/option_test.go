package zerolog

import (
	"testing"

	"github.com/mantzas/patron/http"
	"github.com/mantzas/patron/log"
	"github.com/stretchr/testify/assert"
)

func TestOption_Log(t *testing.T) {
	assert := assert.New(t)
	s, err := http.New("test", []http.Route{http.NewRoute("/", "GET", nil)}, Log(log.DebugLevel))
	assert.NoError(err)
	assert.NotNil(s)
}
