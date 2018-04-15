package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregate(t *testing.T) {
	assert := assert.New(t)
	a := New()
	a.Append(errors.New("Error 1"))
	a.Append(errors.New("Error 2"))
	a.Append(nil)
	a.Append(errors.New("Error 3"))
	assert.Len(a.errors, 3)
	assert.Equal("Error 1\nError 2\nError 3\n", a.Error())
}
