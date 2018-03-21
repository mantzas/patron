package log

import (
	"io"
)

// StdFactory of the std logger
type StdFactory struct {
	w io.Writer
}

// NewStdFactory constructor
func NewStdFactory(w io.Writer) Factory {
	return &StdFactory{w}
}

// Create a std logger
func (sf *StdFactory) Create(f map[string]interface{}) Logger {
	return NewStdLogger(sf.w, f)
}

// CreateSub a std sub logger with defined fields
func (sf *StdFactory) CreateSub(l Logger, f map[string]interface{}) Logger {

	if len(f) == 0 {
		return l
	}

	all := l.Fields()

	for k, v := range f {
		all[k] = v
	}

	return NewStdLogger(sf.w, all)
}
