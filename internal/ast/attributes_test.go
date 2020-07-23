// SPDX-License-Identifier: MIT

package ast

import (
	"net/http"
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
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

func newURIParser(a *assert.Assertion, uri core.URI) (*token.Parser, *messagetest.Result) {
	rslt := messagetest.NewMessageHandler()
	p, err := token.NewParser(rslt.Handler, core.Block{Location: core.Location{URI: uri}})
	a.NotError(err).NotNil(p)
	return p, rslt
}

func TestNumberAttribute(t *testing.T) {
	a := assert.New(t)

	p, rslt := newURIParser(a, "uri1")
	num := &NumberAttribute{}
	attr := &token.Attribute{Value: token.String{Value: "6"}}
	a.NotError(num.DecodeXMLAttr(p, attr))
	a.Equal(num.IntValue(), 6).Equal(0.0, num.FloatValue())
	v, err := num.EncodeXMLAttr()
	a.NotError(err).Equal(v, "6")
	rslt.Handler.Stop()

	p, rslt = newURIParser(a, "uri1")
	num = &NumberAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "6.1"}}
	a.NotError(num.DecodeXMLAttr(p, attr))
	a.Equal(num.IntValue(), 0).Equal(6.1, num.FloatValue())
	v, err = num.EncodeXMLAttr()
	a.NotError(err).Equal(v, "6.1")
	rslt.Handler.Stop()

	p, rslt = newURIParser(a, "uri1")
	num = &NumberAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "6xxy"}}
	a.Error(num.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestBoolAttribute(t *testing.T) {
	a := assert.New(t)

	p, rslt := newURIParser(a, "uri1")
	b := &BoolAttribute{}
	attr := &token.Attribute{Value: token.String{Value: "T"}}
	a.NotError(b.DecodeXMLAttr(p, attr))
	a.True(b.V())
	v, err := b.EncodeXMLAttr()
	a.NotError(err).Equal(v, "true")
	rslt.Handler.Stop()

	p, rslt = newURIParser(a, "uri1")
	b = &BoolAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "xyz"}}
	a.Error(b.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestMethodAttribute(t *testing.T) {
	a := assert.New(t)

	p, rslt := newURIParser(a, "uri1")
	method := &MethodAttribute{}
	attr := &token.Attribute{Value: token.String{Value: http.MethodGet}}
	a.NotError(method.DecodeXMLAttr(p, attr))
	a.Equal(method.V(), http.MethodGet)
	v, err := method.EncodeXMLAttr()
	a.NotError(err).Equal(v, http.MethodGet)
	rslt.Handler.Stop()

	p, rslt = newURIParser(a, "uri1")
	method = &MethodAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "not-exists"}}
	a.Error(method.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestStatusAttribute(t *testing.T) {
	a := assert.New(t)

	p, rslt := newURIParser(a, "uri1")
	status := &StatusAttribute{}
	attr := &token.Attribute{Value: token.String{Value: "201"}}
	a.NotError(status.DecodeXMLAttr(p, attr))
	a.Equal(status.V(), 201)
	v, err := status.EncodeXMLAttr()
	a.NotError(err).Equal(v, "201")
	rslt.Handler.Stop()

	p, rslt = newURIParser(a, "uri1")
	status = &StatusAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "10000"}}
	a.Error(status.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestTypeAttribute(t *testing.T) {
	a := assert.New(t)

	p, rslt := newURIParser(a, "uri1")
	tt := &TypeAttribute{}
	attr := &token.Attribute{Value: token.String{Value: TypeNumber}}
	a.NotError(tt.DecodeXMLAttr(p, attr))
	a.Equal(tt.V(), TypeNumber)
	v, err := tt.EncodeXMLAttr()
	a.NotError(err).Equal(v, TypeNumber)
	rslt.Handler.Stop()

	p, rslt = newURIParser(a, "uri1")
	tt = &TypeAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "10000"}}
	a.Error(tt.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestVersionAttribute(t *testing.T) {
	a := assert.New(t)

	p, rslt := newURIParser(a, "uri1")
	ver := &VersionAttribute{}
	attr := &token.Attribute{Value: token.String{Value: "3.6.1"}}
	a.NotError(ver.DecodeXMLAttr(p, attr))
	a.Equal(ver.V(), "3.6.1")
	v, err := ver.EncodeXMLAttr()
	a.NotError(err).Equal(v, "3.6.1")
	rslt.Handler.Stop()

	p, rslt = newURIParser(a, "uri1")
	ver = &VersionAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "3x"}}
	a.Error(ver.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

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

	p, rslt := newURIParser(a, "uri1")
	now := time.Now().Format(dateFormat)
	date := &DateAttribute{}
	attr := &token.Attribute{
		Name: token.Name{Local: token.String{Value: "n"}},
		Value: token.String{
			Value: now,
			Range: core.Range{End: core.Position{Character: 1}},
		},
	}
	a.NotError(date.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()

	tt, err := date.EncodeXMLAttr()
	a.NotError(err).Equal(tt, now)

	p, rslt = newURIParser(a, "uri1")
	attr.Value.Value = "invalid format"
	err = date.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok := err.(*core.SyntaxError)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1)
	rslt.Handler.Stop()
}

func TestAPIDocVersionAttribute(t *testing.T) {
	a := assert.New(t)

	p, rslt := newURIParser(a, "uri1")
	v := &APIDocVersionAttribute{}
	attr := &token.Attribute{
		Name: token.Name{Local: token.String{Value: "n"}},
		Value: token.String{
			Value: Version,
			Range: core.Range{End: core.Position{Character: 1}},
		},
	}
	a.NotError(v.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()

	vv, err := v.EncodeXMLAttr()
	a.NotError(err).Equal(vv, Version)

	p, rslt = newURIParser(a, "uri1")
	attr.Value.Value = "invalid format"
	err = v.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok := err.(*core.SyntaxError)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1)
	rslt.Handler.Stop()

	// 版本不兼容
	p, rslt = newURIParser(a, "uri1")
	attr.Value.Value = "5.0.0"
	err = v.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1).
		Equal(serr.Err, locale.NewError(locale.ErrInvalidValue))
	rslt.Handler.Stop()
}
