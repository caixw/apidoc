// SPDX-License-Identifier: MIT

package message

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/issue9/assert"
)

var _ error = &SyntaxError{}

func TestSyntaxError_Error(t *testing.T) {
	a := assert.New(t)

	err1 := NewError("file", "", 0, "msg")
	err2 := NewError("file", "field", 0, "msg")
	a.NotEqual(err1.Error(), err2.Error())
}

func TestWithError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("test")
	serr := WithError("file", "field", 1, err)
	a.Equal(serr.Message, err.Error())

	err = &xml.SyntaxError{
		Msg:  "syntaxError",
		Line: 2,
	}
	serr = WithError("file", "field", 1, err)
	a.Equal(serr.Message, "syntaxError").
		Equal(serr.Line, 1+2)
}
