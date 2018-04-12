package zerolog

import (
	"github.com/mantzas/patron/http"
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
)

// Log option for setting the zerolog default logging
func Log(lvl log.Level) http.Option {
	return func(s *http.Service) error {
		err := log.Setup(DefaultFactory(lvl))
		if err != nil {
			return errors.Wrap(err, "failed to set up zerolog default")
		}
		log.Infof("zerolog setup with min level: %d", lvl)
		return nil
	}
}
