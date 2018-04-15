package errors

import (
	"strings"
	"sync"
)

// Aggregate definition of a construct that aggregates multiple errors
// in a safe manner
type Aggregate struct {
	errors []error
	m      sync.Mutex
}

// New creates a new multi error
func New() *Aggregate {
	return &Aggregate{
		errors: []error{},
		m:      sync.Mutex{},
	}
}

// Append a error to the internal collection
func (a *Aggregate) Append(err error) {
	if err == nil {
		return
	}
	a.m.Lock()
	defer a.m.Unlock()
	a.errors = append(a.errors, err)
}

// Count returns the error count
func (a *Aggregate) Count() int {
	a.m.Lock()
	defer a.m.Unlock()
	return len(a.errors)
}

// Error returns the string representation of the errors
// in the internal collection
func (a *Aggregate) Error() string {
	a.m.Lock()
	defer a.m.Unlock()
	b := strings.Builder{}

	for _, err := range a.errors {
		b.WriteString(err.Error())
		b.WriteRune('\n')
	}

	return b.String()
}
