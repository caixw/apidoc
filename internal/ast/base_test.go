// SPDX-License-Identifier: MIT

package ast

import (
	"net/http"
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

func TestNumberAttribute(t *testing.T) {
	a := assert.New(t)
	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

	num := &NumberAttribute{}
	attr := &token.Attribute{Value: token.String{Value: "6"}}
	a.NotError(num.DecodeXMLAttr(p, attr))
	a.Equal(num.IntValue(), 6).Equal(0.0, num.FloatValue())
	v, err := num.EncodeXMLAttr()
	a.NotError(err).Equal(v, "6")

	num = &NumberAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "6.1"}}
	a.NotError(num.DecodeXMLAttr(p, attr))
	a.Equal(num.IntValue(), 0).Equal(6.1, num.FloatValue())
	v, err = num.EncodeXMLAttr()
	a.NotError(err).Equal(v, "6.1")

	num = &NumberAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "6xxy"}}
	a.Error(num.DecodeXMLAttr(p, attr))
}

func TestBoolAttribute(t *testing.T) {
	a := assert.New(t)
	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

	b := &BoolAttribute{}
	attr := &token.Attribute{Value: token.String{Value: "T"}}
	a.NotError(b.DecodeXMLAttr(p, attr))
	a.True(b.V())
	v, err := b.EncodeXMLAttr()
	a.NotError(err).Equal(v, "true")

	b = &BoolAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "xyz"}}
	a.Error(b.DecodeXMLAttr(p, attr))
}

func TestMethodAttribute(t *testing.T) {
	a := assert.New(t)
	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

	method := &MethodAttribute{}
	attr := &token.Attribute{Value: token.String{Value: http.MethodGet}}
	a.NotError(method.DecodeXMLAttr(p, attr))
	a.Equal(method.V(), http.MethodGet)
	v, err := method.EncodeXMLAttr()
	a.NotError(err).Equal(v, http.MethodGet)

	method = &MethodAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "not-exists"}}
	a.Error(method.DecodeXMLAttr(p, attr))
}

func TestStatusAttribute(t *testing.T) {
	a := assert.New(t)
	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

	status := &StatusAttribute{}
	attr := &token.Attribute{Value: token.String{Value: "201"}}
	a.NotError(status.DecodeXMLAttr(p, attr))
	a.Equal(status.V(), 201)
	v, err := status.EncodeXMLAttr()
	a.NotError(err).Equal(v, "201")

	status = &StatusAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "10000"}}
	a.Error(status.DecodeXMLAttr(p, attr))
}

func TestTypeAttribute(t *testing.T) {
	a := assert.New(t)
	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

	tt := &TypeAttribute{}
	attr := &token.Attribute{Value: token.String{Value: TypeNumber}}
	a.NotError(tt.DecodeXMLAttr(p, attr))
	a.Equal(tt.V(), TypeNumber)
	v, err := tt.EncodeXMLAttr()
	a.NotError(err).Equal(v, TypeNumber)

	tt = &TypeAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "10000"}}
	a.Error(tt.DecodeXMLAttr(p, attr))
}

func TestVersionAttribute(t *testing.T) {
	a := assert.New(t)
	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

	ver := &VersionAttribute{}
	attr := &token.Attribute{Value: token.String{Value: "3.6.1"}}
	a.NotError(ver.DecodeXMLAttr(p, attr))
	a.Equal(ver.V(), "3.6.1")
	v, err := ver.EncodeXMLAttr()
	a.NotError(err).Equal(v, "3.6.1")

	ver = &VersionAttribute{}
	attr = &token.Attribute{Value: token.String{Value: "3x"}}
	a.Error(ver.DecodeXMLAttr(p, attr))
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

	p, err := token.NewParser(core.Block{Location: core.Location{URI: core.URI("uri1")}})
	a.NotError(err).NotNil(p)

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
		Name: token.Name{Local: token.String{Value: "n"}},
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

func TestTrimLeftSpace(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		input, output string
	}{
		{},
		{
			input:  `abc`,
			output: `abc`,
		},
		{
			input:  `  abc`,
			output: `abc`,
		},
		{
			input:  "  abc\n",
			output: "abc\n",
		},
		{ // 缩进一个空格
			input:  "  abc\n abc\n",
			output: " abc\nabc\n",
		},
		{ // 缩进一个空格
			input:  "\n  abc\n abc\n",
			output: "\n abc\nabc\n",
		},
		{ // 缩进格式不相同，不会有缩进
			input:  "\t  abc\n abc\n",
			output: "\t  abc\n abc\n",
		},

		{
			input:  "\t  abc\n\t abc\n\t xx\n",
			output: " abc\nabc\nxx\n",
		},
		{
			input:  "\t  abc\n\t abc\nxx\n",
			output: "\t  abc\n\t abc\nxx\n",
		},

		{ // 包含相同的 \t  内容
			input:  "\t  abc\n\t  abc\n\t  xx\n",
			output: "abc\nabc\nxx\n",
		},

		{ // 部分空格相同
			input:  "\t\t  abc\n\t  abc\n\t  xx\n",
			output: "\t  abc\n  abc\n  xx\n",
		},
	}

	for i, item := range data {
		output := trimLeftSpace(item.input)
		a.Equal(output, item.output, "not equal @ %d\nv1=%#v\nv2=%#v\n", i, output, item.output)
	}
}
