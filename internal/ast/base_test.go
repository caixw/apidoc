// SPDX-License-Identifier: MIT

package ast

import (
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

var (
	_ token.AttrDecoder = &Attribute{}
	_ token.AttrDecoder = &NumberAttribute{}
	_ token.AttrDecoder = &BoolAttribute{}
	_ token.AttrDecoder = &MethodAttribute{}
	_ token.AttrDecoder = &StatusAttribute{}
	_ token.AttrDecoder = &TypeAttribute{}
	_ token.AttrDecoder = &DateAttribute{}
	_ token.AttrDecoder = &VersionAttribute{}
	_ token.AttrDecoder = &APIDocVersionAttribute{}

	_ token.AttrEncoder = &Attribute{}
	_ token.AttrEncoder = &NumberAttribute{}
	_ token.AttrEncoder = &BoolAttribute{}
	_ token.AttrEncoder = &MethodAttribute{}
	_ token.AttrEncoder = &StatusAttribute{}
	_ token.AttrEncoder = &TypeAttribute{}
	_ token.AttrEncoder = &DateAttribute{}
	_ token.AttrEncoder = &VersionAttribute{}
	_ token.AttrEncoder = &APIDocVersionAttribute{}
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

func TestDateAttribute(t *testing.T) {
	a := assert.New(t)

	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

	now := time.Now().Format(dateFormat)
	date := &DateAttribute{}
	attr := &token.Attribute{
		Name: token.String{Value: "n"},
		Value: token.String{
			Value: now,
			Range: core.Range{End: core.Position{Character: 1}},
		},
	}
	a.NotError(date.DecodeXMLAttr(p, attr))

	tt, err := date.EncodeXMLAttr()
	a.NotError(err).Equal(tt, now)

	attr.Value.Value = "invalid format"
	err = date.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok := err.(*core.SyntaxError)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1)
}

func TestAPIDocVersionAttribute(t *testing.T) {
	a := assert.New(t)

	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

	v := &APIDocVersionAttribute{}
	attr := &token.Attribute{
		Name: token.String{Value: "n"},
		Value: token.String{
			Value: Version,
			Range: core.Range{End: core.Position{Character: 1}},
		},
	}
	a.NotError(v.DecodeXMLAttr(p, attr))

	vv, err := v.EncodeXMLAttr()
	a.NotError(err).Equal(vv, Version)

	attr.Value.Value = "invalid format"
	err = v.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok := err.(*core.SyntaxError)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1)

	// 版本不兼容
	attr.Value.Value = "5.0.0"
	err = v.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1).
		Equal(serr.Err, locale.NewError(locale.ErrInvalidValue))
}
