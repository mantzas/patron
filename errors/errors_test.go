package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("TEST")
	assert.EqualError(t, err, "TEST")
}

func TestErrorf(t *testing.T) {
	err := Errorf("TEST %s", "one")
	assert.EqualError(t, err, "TEST one")
}

func TestWrap(t *testing.T) {
	err := Wrap(New("TEST"), "Wrap")
	assert.EqualError(t, err, "Wrap: TEST")
}

func TestWrapf(t *testing.T) {
	err := Wrapf(New("TEST"), "Wrap %s", "error")
	assert.EqualError(t, err, "Wrap error: TEST")
}
