package errors

import (
	"strings"
	"sync"
)

// Aggregate for aggregating errors into one.
// The aggregation is goroutine safe.
type Aggregate struct {
	sync.Mutex
	errors []error
}

// NewAggregate creates a new aggregate error.
func NewAggregate() *Aggregate {
	return &Aggregate{errors: []error{}}
}

// Append a error to the internal collection.
func (a *Aggregate) Append(err error) {
	if err == nil {
		return
	}
	a.Lock()
	defer a.Unlock()
	a.errors = append(a.errors, err)
}

// Count returns the count of the aggregated errors.
func (a *Aggregate) Count() int {
	a.Lock()
	defer a.Unlock()
	return len(a.errors)
}

// Error returns the string representation of the aggregated errors.
func (a *Aggregate) Error() string {
	a.Lock()
	defer a.Unlock()
	b := strings.Builder{}
	for _, err := range a.errors {
		b.WriteString(err.Error())
		b.WriteRune('\n')
	}
	return b.String()
}
