// SPDX-License-Identifier: MIT

package ast

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/internal/token"
)

var (
	_ token.AttrDecoder = &Attribute{}
	_ token.AttrDecoder = &NumberAttribute{}
	_ token.AttrDecoder = &BoolAttribute{}

	_ token.AttrEncoder = &Attribute{}
	_ token.AttrEncoder = &NumberAttribute{}
	_ token.AttrEncoder = &BoolAttribute{}

	_ token.Decoder = &Element{}
	_ token.Encoder = &Element{}
)

func TestIsValidMethod(t *testing.T) {
	a := assert.New(t)

	a.True(isValidMethod("GET"))
	a.True(isValidMethod("get"))
	a.False(isValidMethod("not-exists"))
}

func TestIsValidStatus(t *testing.T) {
	a := assert.New(t)

	a.True(isValidStatus(100))
	a.True(isValidStatus(500))
	a.False(isValidStatus(1000))
}
