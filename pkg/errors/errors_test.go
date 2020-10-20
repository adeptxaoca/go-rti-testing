package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetType(t *testing.T) {
	err := BadRequest.New("error")
	assert.Equal(t, BadRequest, GetType(err))
	err = New("error")
	assert.Equal(t, Internal, GetType(err))
}

func TestWrap(t *testing.T) {
	err := New("some kind of error")
	err = Wrap(err, "comment on error")
	assert.Equal(t, Internal, GetType(err))
	assert.Equal(t, "comment on error: some kind of error", err.Error())
}

func TestCause(t *testing.T) {
	err := BadRequest.New("original")
	err = Wrap(err, "second")
	err = Wrap(err, "another")
	original := Cause(err)
	require.NotNil(t, original)
	assert.Equal(t, "original", original.Error())

	err = MethodNotAllowed.New("one")
	original = Cause(err)
	require.NotNil(t, original)
	assert.Equal(t, "one", original.Error())
}
