package zerolog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mantzas/patron/log"
	"github.com/rs/zerolog"
)

// Factory implementation of zerolog.
type Factory struct {
	logger *zerolog.Logger
	lvl    log.Level
}

// NewFactory creates a new zerolog factory.
func NewFactory(l *zerolog.Logger, lvl log.Level) log.Factory {
	return &Factory{logger: l, lvl: lvl}
}

// DefaultFactory creates a zerolog factory with default settings.
func DefaultFactory(lvl log.Level) log.Factory {
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zl := zerolog.New(os.Stdout).With().Timestamp().Logger().Hook(sourceHook{})
	return NewFactory(&zl, lvl)
}

// Create a new logger.
func (zf *Factory) Create(f map[string]interface{}) log.Logger {
	return NewLogger(zf.logger, zf.lvl, f)
}

type sourceHook struct{}

func (sh sourceHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	k, v, ok := sourceFields()
	if !ok {
		return
	}
	e.Str(k, v)
}

func sourceFields() (key string, src string, ok bool) {
	_, file, line, ok := runtime.Caller(7)
	if !ok {
		return
	}
	src = getSource(file, line)
	key = "src"
	ok = true
	return
}

func getSource(file string, line int) (src string) {
	if file == "" {
		return
	}
	d, f := filepath.Split(file)
	d = path.Base(d)
	if d == "." || d == "" {
		src = fmt.Sprintf("%s:%d", f, line)
	} else {
		src = fmt.Sprintf("%s/%s:%d", d, f, line)
	}
	return
}
