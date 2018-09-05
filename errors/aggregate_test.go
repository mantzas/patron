package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregate(t *testing.T) {
	a := Aggregate(New("Error 1"), New("Error 2"), nil, New("Error 3"))
	assert.Len(t, a, 3)
	assert.Equal(t, "Error 1\nError 2\nError 3\n", a.Error())
}
