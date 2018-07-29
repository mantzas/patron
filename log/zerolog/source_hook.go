package zerolog

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
)

type sourceHook struct{}

func (h sourceHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {

	if e == nil {
		return
	}

	_, file, line, ok := runtime.Caller(5)
	if !ok {
		return
	}

	src := getSource(file, line)
	if src == "" {
		return
	}

	e.Fields(map[string]interface{}{"src": src})
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
