// SPDX-License-Identifier: MIT

package core

import (
	"errors"
	"testing"

	"github.com/issue9/assert"
)

var _ error = &SyntaxError{}

func TestNewLocaleError(t *testing.T) {
	a := assert.New(t)

	err1 := NewSyntaxError(Location{}, "", "msg")
	err2 := NewSyntaxError(Location{}, "field", "msg")
	a.NotEqual(err1.Error(), err2.Error())
}

func TestWithError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("test")
	serr := NewSyntaxErrorWithError(Location{}, "field", err)
	a.Equal(serr.Err.Error(), err.Error())
}
