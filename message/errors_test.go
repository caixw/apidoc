// SPDX-License-Identifier: MIT

package message

import (
	"errors"
	"testing"

	"github.com/issue9/assert"
)

var _ error = &SyntaxError{}

func TestNewLocaleError(t *testing.T) {
	a := assert.New(t)

	err1 := NewLocaleError("file", "", 0, "msg")
	err2 := NewLocaleError("file", "field", 0, "msg")
	a.NotEqual(err1.Error(), err2.Error())
}

func TestWithError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("test")
	serr := WithError("file", "field", 1, err)
	a.Equal(serr.Message, err.Error())
}
