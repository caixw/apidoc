// SPDX-License-Identifier: MIT

package core

import (
	"errors"
	"os"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/internal/locale"
)

var _ error = &Error{}

func TestError(t *testing.T) {
	a := assert.New(t, false)

	err1 := NewError("msg")
	err2 := NewError("msg").WithField("field")
	a.NotEqual(err1.Error(), err2.Error())
}

func TestWithError(t *testing.T) {
	a := assert.New(t, false)

	err := errors.New("test")
	serr := WithError(err).WithField("field")
	a.Equal(serr.Err, err)

	serr2 := WithError(serr).WithLocation(Location{URI: "uri"})
	a.Equal(serr2.Err, err)
}

func TestError_AddTypes(t *testing.T) {
	a := assert.New(t, false)
	loc := Location{}

	err := loc.WithError(errors.New("err1"))
	err.AddTypes(ErrorTypeDeprecated)
	a.Equal(err.Types, []ErrorType{ErrorTypeDeprecated})
	err.AddTypes(ErrorTypeDeprecated)
	a.Equal(err.Types, []ErrorType{ErrorTypeDeprecated})

	err.AddTypes(ErrorTypeUnused)
	a.Equal(err.Types, []ErrorType{ErrorTypeDeprecated, ErrorTypeUnused})
}

func TestError_Is_Unwrap(t *testing.T) {
	a := assert.New(t, false)

	err := WithError(os.ErrExist).WithField("field")
	a.True(errors.Is(err, os.ErrExist))

	a.Equal(errors.Unwrap(err), os.ErrExist)
}

func TestError_Relate(t *testing.T) {
	a := assert.New(t, false)

	err := NewError(locale.ErrInvalidUTF8Character)
	a.Empty(err.Related)
	err.Relate(Location{}, "msg")
	a.Equal(1, len(err.Related)).Equal(err.Related[0].Message, "msg")

	err.Relate(Location{}, "msg2")
	a.Equal(2, len(err.Related)).Equal(err.Related[1].Message, "msg2")
}
