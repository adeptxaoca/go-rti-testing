package errors

// Package errors provides simple error handling primitives.
// Package based on ideas https://github.com/henrmota/errors-handling-example

import (
	"errors"
	"fmt"

	gErrors "github.com/pkg/errors"
)

type ErrorType uint

const (
	Internal ErrorType = iota
	BadRequest
	UnsupportedMediaType
	MethodNotAllowed
)

const (
	MsgUnsupportedMediaType = "content-type header is not '%s'"
)

type Error struct {
	errType       ErrorType
	originalError error
	context       errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// Error returns the message of a mpError
func (error Error) Error() string {
	return error.originalError.Error()
}

// New creates a new mpError
func (errType ErrorType) New(msg string) error {
	return Error{errType: errType, originalError: gErrors.New(msg)}
}

// New creates a new mpError with formatted message
func (errType ErrorType) Newf(msg string, args ...interface{}) error {
	return Error{errType: errType, originalError: fmt.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (errType ErrorType) Wrap(err error, msg string) error {
	return errType.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (errType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return Error{errType: errType, originalError: gErrors.Wrapf(err, msg, args...)}
}

// New returns an error that formats as the given text.
func New(msg string) error {
	return errors.New(msg)
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	for err != nil {
		mrpErr, ok := err.(Error)
		if !ok {
			break
		}
		err = gErrors.Cause(mrpErr.originalError)
	}
	return err
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := gErrors.Wrapf(err, msg, args...)
	if mrpErr, ok := err.(Error); ok {
		return Error{
			errType:       mrpErr.errType,
			originalError: wrappedError,
			context:       mrpErr.context,
		}
	}

	return Error{errType: Internal, originalError: wrappedError}
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if mpErr, ok := err.(Error); ok {
		return mpErr.errType
	}
	return Internal
}
