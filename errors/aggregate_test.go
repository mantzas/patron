package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregate(t *testing.T) {
	assert := assert.New(t)
	a := NewAggregate()
	a.Append(New("Error 1"))
	a.Append(New("Error 2"))
	a.Append(nil)
	a.Append(New("Error 3"))
	assert.Equal(a.Count(), 3)
	assert.Equal("Error 1\nError 2\nError 3\n", a.Error())
}
