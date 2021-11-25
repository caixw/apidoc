// SPDX-License-Identifier: MIT

package ast

import (
	"net/http"
	"testing"
	"time"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

var (
	_ xmlenc.AttrDecoder = &Attribute{}
	_ xmlenc.AttrDecoder = &NumberAttribute{}
	_ xmlenc.AttrDecoder = &BoolAttribute{}
	_ xmlenc.AttrDecoder = &MethodAttribute{}
	_ xmlenc.AttrDecoder = &StatusAttribute{}
	_ xmlenc.AttrDecoder = &TypeAttribute{}
	_ xmlenc.AttrDecoder = &DateAttribute{}
	_ xmlenc.AttrDecoder = &VersionAttribute{}
	_ xmlenc.AttrDecoder = &APIDocVersionAttribute{}

	_ xmlenc.AttrEncoder = &Attribute{}
	_ xmlenc.AttrEncoder = &NumberAttribute{}
	_ xmlenc.AttrEncoder = &BoolAttribute{}
	_ xmlenc.AttrEncoder = &MethodAttribute{}
	_ xmlenc.AttrEncoder = &StatusAttribute{}
	_ xmlenc.AttrEncoder = &TypeAttribute{}
	_ xmlenc.AttrEncoder = &DateAttribute{}
	_ xmlenc.AttrEncoder = &VersionAttribute{}
	_ xmlenc.AttrEncoder = &APIDocVersionAttribute{}

	_ core.Searcher = Attribute{}
)

func TestNumberAttribute(t *testing.T) {
	a := assert.New(t, false)

	p, rslt := newParser(a, "", "uri1")
	num := &NumberAttribute{}
	attr := &xmlenc.Attribute{Value: xmlenc.String{Value: "6"}}
	a.NotError(num.DecodeXMLAttr(p, attr))
	a.Equal(num.IntValue(), 6).Equal(0.0, num.FloatValue())
	v, err := num.EncodeXMLAttr()
	a.NotError(err).Equal(v, "6")
	rslt.Handler.Stop()

	p, rslt = newParser(a, "", "uri1")
	num = &NumberAttribute{}
	attr = &xmlenc.Attribute{Value: xmlenc.String{Value: "6.1"}}
	a.NotError(num.DecodeXMLAttr(p, attr))
	a.Equal(num.IntValue(), 0).Equal(6.1, num.FloatValue())
	v, err = num.EncodeXMLAttr()
	a.NotError(err).Equal(v, "6.1")
	rslt.Handler.Stop()

	p, rslt = newParser(a, "", "uri1")
	num = &NumberAttribute{}
	attr = &xmlenc.Attribute{Value: xmlenc.String{Value: "6xxy"}}
	a.Error(num.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestBoolAttribute(t *testing.T) {
	a := assert.New(t, false)

	p, rslt := newParser(a, "", "uri1")
	b := &BoolAttribute{}
	attr := &xmlenc.Attribute{Value: xmlenc.String{Value: "T"}}
	a.NotError(b.DecodeXMLAttr(p, attr))
	a.True(b.V())
	v, err := b.EncodeXMLAttr()
	a.NotError(err).Equal(v, "true")
	rslt.Handler.Stop()

	p, rslt = newParser(a, "", "uri1")
	b = &BoolAttribute{}
	attr = &xmlenc.Attribute{Value: xmlenc.String{Value: "xyz"}}
	a.Error(b.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestMethodAttribute(t *testing.T) {
	a := assert.New(t, false)

	p, rslt := newParser(a, "", "uri1")
	method := &MethodAttribute{}
	attr := &xmlenc.Attribute{Value: xmlenc.String{Value: http.MethodGet}}
	a.NotError(method.DecodeXMLAttr(p, attr))
	a.Equal(method.V(), http.MethodGet)
	v, err := method.EncodeXMLAttr()
	a.NotError(err).Equal(v, http.MethodGet)
	rslt.Handler.Stop()

	p, rslt = newParser(a, "", "uri1")
	method = &MethodAttribute{}
	attr = &xmlenc.Attribute{Value: xmlenc.String{Value: "not-exists"}}
	a.Error(method.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestStatusAttribute(t *testing.T) {
	a := assert.New(t, false)

	p, rslt := newParser(a, "", "uri1")
	status := &StatusAttribute{}
	attr := &xmlenc.Attribute{Value: xmlenc.String{Value: "201"}}
	a.NotError(status.DecodeXMLAttr(p, attr))
	a.Equal(status.V(), 201)
	v, err := status.EncodeXMLAttr()
	a.NotError(err).Equal(v, "201")
	rslt.Handler.Stop()

	p, rslt = newParser(a, "", "uri1")
	status = &StatusAttribute{}
	attr = &xmlenc.Attribute{Value: xmlenc.String{Value: "10000"}}
	a.Error(status.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestTypeAttribute(t *testing.T) {
	a := assert.New(t, false)

	p, rslt := newParser(a, "", "uri1")
	tt := &TypeAttribute{}
	attr := &xmlenc.Attribute{Value: xmlenc.String{Value: TypeNumber}}
	a.NotError(tt.DecodeXMLAttr(p, attr))
	a.Equal(tt.V(), TypeNumber)
	v, err := tt.EncodeXMLAttr()
	a.NotError(err).Equal(v, TypeNumber)
	rslt.Handler.Stop()

	p, rslt = newParser(a, "", "uri1")
	tt = &TypeAttribute{}
	attr = &xmlenc.Attribute{Value: xmlenc.String{Value: "10000"}}
	a.Error(tt.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestVersionAttribute(t *testing.T) {
	a := assert.New(t, false)

	p, rslt := newParser(a, "", "uri1")
	ver := &VersionAttribute{}
	attr := &xmlenc.Attribute{Value: xmlenc.String{Value: "3.6.1"}}
	a.NotError(ver.DecodeXMLAttr(p, attr))
	a.Equal(ver.V(), "3.6.1")
	v, err := ver.EncodeXMLAttr()
	a.NotError(err).Equal(v, "3.6.1")
	rslt.Handler.Stop()

	p, rslt = newParser(a, "", "uri1")
	ver = &VersionAttribute{}
	attr = &xmlenc.Attribute{Value: xmlenc.String{Value: "3x"}}
	a.Error(ver.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()
}

func TestIsValidMethod(t *testing.T) {
	a := assert.New(t, false)

	a.True(isValidMethod("GET"))
	a.False(isValidMethod("not-exists"))
}

func TestIsValidStatus(t *testing.T) {
	a := assert.New(t, false)

	a.True(isValidStatus(100))
	a.True(isValidStatus(500))
	a.False(isValidStatus(1000))
}

func TestDateAttribute(t *testing.T) {
	a := assert.New(t, false)

	p, rslt := newParser(a, "", "uri1")
	now := time.Now().Format(dateFormat)
	date := &DateAttribute{}
	attr := &xmlenc.Attribute{
		Name: xmlenc.Name{Local: xmlenc.String{Value: "n"}},
		Value: xmlenc.String{
			Value: now,
			Location: core.Location{
				URI:   "uri1",
				Range: core.Range{End: core.Position{Character: 1}},
			},
		},
	}
	a.NotError(date.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()

	tt, err := date.EncodeXMLAttr()
	a.NotError(err).Equal(tt, now)

	p, rslt = newParser(a, "", "uri1")
	attr.Value.Value = "invalid format"
	err = date.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok := err.(*core.Error)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1)
	rslt.Handler.Stop()
}

func TestAPIDocVersionAttribute(t *testing.T) {
	a := assert.New(t, false)

	p, rslt := newParser(a, "", "uri1")
	v := &APIDocVersionAttribute{}
	attr := &xmlenc.Attribute{
		Name: xmlenc.Name{Local: xmlenc.String{Value: "n"}},
		Value: xmlenc.String{
			Value: Version,
			Location: core.Location{
				URI:   "uri1",
				Range: core.Range{End: core.Position{Character: 1}},
			},
		},
	}
	a.NotError(v.DecodeXMLAttr(p, attr))
	rslt.Handler.Stop()

	vv, err := v.EncodeXMLAttr()
	a.NotError(err).Equal(vv, Version)

	p, rslt = newParser(a, "", "uri1")
	attr.Value.Value = "invalid format"
	err = v.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok := err.(*core.Error)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1)
	rslt.Handler.Stop()

	// 版本不兼容
	p, rslt = newParser(a, "", "uri1")
	attr.Value.Value = "5.0.0"
	err = v.DecodeXMLAttr(p, attr)
	a.Error(err)
	serr, ok = err.(*core.Error)
	a.True(ok).
		Equal(serr.Location.URI, "uri1").
		Equal(serr.Location.Range.End.Character, 1).
		Equal(serr.Err, locale.NewError(locale.ErrInvalidValue))
	rslt.Handler.Stop()
}
