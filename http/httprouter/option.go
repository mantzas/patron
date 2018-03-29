package httprouter

import (
	"github.com/mantzas/patron/http"
)

// Handler option for setting the handler generator
func Handler() http.Option {
	return func(s http.Service) error {
		s.HandlerGen = CreateHandler
		return nil
	}
}
