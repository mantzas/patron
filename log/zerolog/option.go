package zerolog

import (
	"github.com/mantzas/patron/http"
	"github.com/mantzas/patron/log"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Log option for setting the zerolog default logging
func Log(l zerolog.Level) http.Option {
	return func(s http.Service) error {
		err := log.Setup(DefaultFactory(l))
		if err != nil {
			return errors.Wrap(err, "failed to set up zerolog default")
		}
		return nil
	}
}
