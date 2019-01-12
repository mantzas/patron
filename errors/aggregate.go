package errors

import (
	"strings"
)

// Aggregate of errors into one.
type aggregate []error

// Error returns the string representation of the aggregated errors.
func (a aggregate) Error() string {
	b := strings.Builder{}
	for _, err := range a {
		b.WriteString(err.Error())
		b.WriteRune('\n')
	}
	return b.String()
}

// Aggregate errors into one error.
func Aggregate(ee ...error) error {
	if len(ee) == 0 {
		return nil
	}
	agr := make(aggregate, len(ee))
	for i := 0; i < len(ee); i++ {
		if ee[i] == nil {
			continue
		}
		agr[i] = ee[i]
	}
	if len(agr) == 0 {
		return nil
	}
	return agr
}
