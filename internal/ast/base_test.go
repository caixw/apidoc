// SPDX-License-Identifier: MIT

package ast

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/internal/token"
)

var (
	_ token.AttrDecoder = &Attribute{}
	_ token.AttrDecoder = &NumberAttribute{}
	_ token.AttrDecoder = &BoolAttribute{}
	_ token.AttrDecoder = &MethodAttribute{}
	_ token.AttrDecoder = &StatusAttribute{}
	_ token.AttrDecoder = &TypeAttribute{}
	_ token.AttrDecoder = &VersionAttribute{}
	_ token.AttrDecoder = &APIDocVersionAttribute{}

	_ token.AttrEncoder = &Attribute{}
	_ token.AttrEncoder = &NumberAttribute{}
	_ token.AttrEncoder = &BoolAttribute{}
	_ token.AttrEncoder = &MethodAttribute{}
	_ token.AttrEncoder = &StatusAttribute{}
	_ token.AttrEncoder = &TypeAttribute{}
	_ token.AttrEncoder = &VersionAttribute{}
	_ token.AttrEncoder = &APIDocVersionAttribute{}

	_ token.Decoder = &Element{}
	_ token.Encoder = &Element{}
)

func TestIsValidMethod(t *testing.T) {
	a := assert.New(t)

	a.True(isValidMethod("GET"))
	a.False(isValidMethod("not-exists"))
}

func TestIsValidStatus(t *testing.T) {
	a := assert.New(t)

	a.True(isValidStatus(100))
	a.True(isValidStatus(500))
	a.False(isValidStatus(1000))
}
