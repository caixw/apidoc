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

	err1 := NewLocaleError(Location{}, "", "msg")
	err2 := NewLocaleError(Location{}, "field", "msg")
	a.NotEqual(err1.Error(), err2.Error())
}

func TestWithError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("test")
	serr := WithError(Location{}, "field", err)
	a.Equal(serr.Message, err.Error())
}
