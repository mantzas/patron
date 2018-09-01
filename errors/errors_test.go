package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	err := New("TEST")
	assert.EqualError(err, "TEST")
}

func TestErrorf(t *testing.T) {
	assert := assert.New(t)
	err := Errorf("TEST %s", "one")
	assert.EqualError(err, "TEST one")
}

func TestWrap(t *testing.T) {
	assert := assert.New(t)
	err := Wrap(New("TEST"), "Wrap")
	assert.EqualError(err, "Wrap: TEST")
}

func TestWrapf(t *testing.T) {
	assert := assert.New(t)
	err := Wrapf(New("TEST"), "Wrap %s", "error")
	assert.EqualError(err, "Wrap error: TEST")
}
