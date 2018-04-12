package http

import (
	"testing"

	"github.com/mantzas/patron/log"
	"github.com/mantzas/patron/log/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestPorts(t *testing.T) {
	log.Setup(zerolog.DefaultFactory(log.DebugLevel))
	assert := assert.New(t)
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"success", 30000, false},
		{"error for port number out of range", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, err := New("test", getRoutes("/3"), Ports(tt.port))

			if tt.wantErr {
				assert.Nil(s)
				assert.Error(err)
			} else {
				assert.NotNil(s)
				assert.NoError(err)
			}
		})
	}
}
