// SPDX-License-Identifier: MIT

package core

import (
	"errors"
	"os"
	"testing"

	"github.com/issue9/assert"
)

var _ error = &Error{}

func TestError(t *testing.T) {
	a := assert.New(t)

	err1 := NewError("msg")
	err2 := NewError("msg").WithField("field")
	a.NotEqual(err1.Error(), err2.Error())
}

func TestWithError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("test")
	serr := WithError(err).WithField("field")
	a.Equal(serr.Err, err)

	serr2 := WithError(serr).WithLocation(Location{URI: "uri"})
	a.Equal(serr2.Err, err)
}

func TestError_AddTypes(t *testing.T) {
	a := assert.New(t)
	loc := Location{}

	err := loc.WithError(errors.New("err1"))
	err.AddTypes(ErrorTypeDeprecated)
	a.Equal(err.Types, []ErrorType{ErrorTypeDeprecated})
	err.AddTypes(ErrorTypeDeprecated)
	a.Equal(err.Types, []ErrorType{ErrorTypeDeprecated})

	err.AddTypes(ErrorTypeUnused)
	a.Equal(err.Types, []ErrorType{ErrorTypeDeprecated, ErrorTypeUnused})
}

func TestErrorr_Is_Unwrap(t *testing.T) {
	a := assert.New(t)

	err := WithError(os.ErrExist).WithField("field")
	a.True(errors.Is(err, os.ErrExist))

	a.Equal(errors.Unwrap(err), os.ErrExist)
}
