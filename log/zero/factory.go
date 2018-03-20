package zero

import (
	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
)

// Factory of the zero logger
type Factory struct {
	l *zerolog.Logger
}

// NewFactory returns a new zero logger factory
func NewFactory(l *zerolog.Logger) log.Factory {
	return &Factory{l}
}

// Create a zero logger
func (f *Factory) Create() log.Logger {
	return NewLogger(f.l, make(map[string]interface{}))
}

// CreateWithFields a zero logger with defined fields
func (f *Factory) CreateWithFields(fields map[string]interface{}) log.Logger {

	if len(fields) == 0 {
		return NewLogger(f.l, make(map[string]interface{}))
	}

	return NewLogger(f.createZerologger(fields), fields)
}

// CreateSub a zero sub logger with defined fields
func (f *Factory) CreateSub(logger log.Logger, fields map[string]interface{}) log.Logger {

	if len(fields) == 0 {
		return logger
	}

	all := logger.Fields()

	for k, v := range fields {
		all[k] = v
	}

	return NewLogger(f.createZerologger(all), all)
}

func (f *Factory) createZerologger(fields map[string]interface{}) *zerolog.Logger {

	l := f.l.With().Fields(fields).Logger()
	return &l
}
