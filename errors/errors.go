package errors

import (
	"github.com/pkg/errors"
)

// New creates a new error.
func New(msg string) error {
	return errors.New(msg)
}

// Errorf returns a error with a formated message.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// Wrap returns a error that wraps a error and augments it with a new message.
func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

// Wrapf returns a error that wraps a error and augments it with a new formatted message.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
