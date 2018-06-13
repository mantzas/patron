package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_metricRoute(t *testing.T) {
	assert := assert.New(t)
	route := metricRoute()
	assert.Equal(http.MethodGet, route.Method)
	assert.Equal("/metric", route.Pattern)
	assert.NotNil(route.Handler)
	assert.False(route.Trace)
}
