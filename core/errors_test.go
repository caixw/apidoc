// SPDX-License-Identifier: MIT

package core

import (
	"errors"
	"os"
	"testing"

	"github.com/issue9/assert"
)

var _ error = &SyntaxError{}

func TestNewSyntaxError(t *testing.T) {
	a := assert.New(t)

	err1 := NewSyntaxError(Location{}, "", "msg")
	err2 := NewSyntaxError(Location{}, "field", "msg")
	a.NotEqual(err1.Error(), err2.Error())
}

func TestNewSyntaxErrorWithError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("test")
	serr := NewSyntaxErrorWithError(Location{}, "field", err)
	a.Equal(serr.Err, err)

	serr2 := NewSyntaxErrorWithError(Location{}, "", serr)
	a.Equal(serr2.Err, err)
}

func TestSyntaxErrorr_Is_Unwrap(t *testing.T) {
	a := assert.New(t)

	err := NewSyntaxErrorWithError(Location{}, "field", os.ErrExist)
	a.True(errors.Is(err, os.ErrExist))

	a.Equal(errors.Unwrap(err), os.ErrExist)
}
