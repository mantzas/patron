package errors

import (
	"strings"
	"sync"

	"github.com/mantzas/patron/log"
)

// Aggregate for aggregating errors into one.
// The aggregation is goroutine safe.
type Aggregate struct {
	sync.Mutex
	errors []error
}

// New creates a new aggregate error.
func New() *Aggregate {
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
		_, err1 := b.WriteString(err.Error())
		if err1 != nil {
			log.Errorf("failed to write %v with error %v", err, err1)
		}
		_, err1 = b.WriteRune('\n')
		if err1 != nil {
			log.Errorf("failed to write newline")
		}
	}
	return b.String()
}
