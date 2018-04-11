package httprouter

import (
	"github.com/mantzas/patron/http"
	"github.com/mantzas/patron/log"
)

// Handler option for setting the handler generator
func Handler() http.Option {
	return func(s *http.Service) error {
		s.HandlerGen = CreateHandler
		log.Info("httprouter handler set up")
		return nil
	}
}
