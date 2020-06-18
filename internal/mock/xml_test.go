// SPDX-License-Identifier: MIT

package mock

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/token"
)

func TestValidXML(t *testing.T) {
	a := assert.New(t)

	for _, item := range data {
		err := validXML(item.XMLNS, item.Type, []byte(item.XML))
		a.NotError(err, "测试 %s 时返回错误 %s", item.Title, err)
	}

	p := &ast.Request{
		Name: &ast.Attribute{Value: token.String{Value: "root"}},
		Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}},
		Items: []*ast.Param{
			{
				Name: &ast.Attribute{Value: token.String{Value: "id"}},
				Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
				XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
			},
			{
				Name: &ast.Attribute{Value: token.String{Value: "desc"}},
				Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
				XML:  ast.XML{XMLExtract: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
			},
		},
	}
	content := `<root id="1024"><desc>1024</desc></root>`
	a.Error(validXML(nil, p, []byte(content)))
}

func TestBuildXML(t *testing.T) {
	a := assert.New(t)

	for _, item := range data {
		data, err := buildXML(item.XMLNS, item.Type, indent, testOptions)
		a.NotError(err, "测试 %s 返回了错误信息 %s", item.Title, err).
			Equal(string(data), item.XML, "测试 %s 返回的数据不相等\nv1:%s\nv2:%s\n", item.Title, string(data), item.XML)
	}
}

func TestValidXMLName(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		name       xml.Name
		ns         []*ast.XMLNamespace
		param      *ast.Param
		allowArray bool
	}{
		{
			name: xml.Name{Local: "n1"},
			param: &ast.Param{
				Name: &ast.Attribute{Value: token.String{Value: "n1"}},
			},
		},
		{
			name: xml.Name{Local: "n1"},
			param: &ast.Param{
				Name: &ast.Attribute{Value: token.String{Value: "n1"}},
			},
			allowArray: true,
		},

		{
			name: xml.Name{Local: "n1"},
			param: &ast.Param{
				Name: &ast.Attribute{Value: token.String{Value: "n1"}},
				XML: ast.XML{
					XMLWrapped: &ast.Attribute{Value: token.String{Value: "parent"}},
				},
			},
		},
		{
			name: xml.Name{Local: "n1"},
			param: &ast.Param{
				Name: &ast.Attribute{Value: token.String{Value: "n1"}},
				XML: ast.XML{
					XMLWrapped: &ast.Attribute{Value: token.String{Value: "parent"}},
				},
			},
			allowArray: true,
		},

		// wrapped=parent
		{
			name: xml.Name{Local: "parent", Space: "urn"},
			ns: []*ast.XMLNamespace{
				{
					URN:    &ast.Attribute{Value: token.String{Value: "urn"}},
					Prefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			param: &ast.Param{
				Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
				Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				XML: ast.XML{
					XMLWrapped:  &ast.Attribute{Value: token.String{Value: "parent"}},
					XMLNSPrefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			allowArray: true,
		},
		{
			name: xml.Name{Local: "n1", Space: "urn"},
			ns: []*ast.XMLNamespace{
				{
					URN:    &ast.Attribute{Value: token.String{Value: "urn"}},
					Prefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			param: &ast.Param{
				Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
				Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				XML: ast.XML{
					XMLWrapped:  &ast.Attribute{Value: token.String{Value: "parent"}},
					XMLNSPrefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
		},

		// wrapped=parent>n2
		{
			name: xml.Name{Local: "parent", Space: "urn"},
			ns: []*ast.XMLNamespace{
				{
					URN:    &ast.Attribute{Value: token.String{Value: "urn"}},
					Prefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			param: &ast.Param{
				Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
				Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				XML: ast.XML{
					XMLWrapped:  &ast.Attribute{Value: token.String{Value: "parent>n2"}},
					XMLNSPrefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			allowArray: true,
		},
		{
			name: xml.Name{Local: "n2", Space: "urn"},
			ns: []*ast.XMLNamespace{
				{
					URN:    &ast.Attribute{Value: token.String{Value: "urn"}},
					Prefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			param: &ast.Param{
				Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
				Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				XML: ast.XML{
					XMLWrapped:  &ast.Attribute{Value: token.String{Value: "parent>n2"}},
					XMLNSPrefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
		},

		// wrapped=>n2
		{
			name: xml.Name{Local: "", Space: ""},
			ns: []*ast.XMLNamespace{
				{
					URN:    &ast.Attribute{Value: token.String{Value: "urn"}},
					Prefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			param: &ast.Param{
				Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
				Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				XML: ast.XML{
					XMLWrapped:  &ast.Attribute{Value: token.String{Value: ">n2"}},
					XMLNSPrefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			allowArray: true,
		},
		{
			name: xml.Name{Local: "n2", Space: "urn"},
			ns: []*ast.XMLNamespace{
				{
					URN:    &ast.Attribute{Value: token.String{Value: "urn"}},
					Prefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
			param: &ast.Param{
				Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
				Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				XML: ast.XML{
					XMLWrapped:  &ast.Attribute{Value: token.String{Value: ">n2"}},
					XMLNSPrefix: &ast.Attribute{Value: token.String{Value: "p"}},
				},
			},
		},
	}

	for i, item := range data {
		validator := &xmlValidator{namespaces: item.ns}
		a.True(validator.validXMLName(item.name, item.param, item.allowArray), "false at %d", i)
	}
}

func TestParseXMLName(t *testing.T) {
	a := assert.New(t)

	n := parseXMLWrappedName(&ast.Param{
		Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
		Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
	}, false)
	a.Equal(n, "n1")
	n = parseXMLWrappedName(&ast.Param{
		Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
		Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
	}, true)
	a.Empty(n)

	n = parseXMLWrappedName(&ast.Param{
		Name: &ast.Attribute{Value: token.String{Value: "n1"}},
		XML: ast.XML{
			XMLWrapped: &ast.Attribute{Value: token.String{Value: "parent"}},
		},
	}, true)
	a.Equal(n, "n1")

	// wrapped = parent
	n = parseXMLWrappedName(&ast.Param{
		Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
		Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		XML: ast.XML{
			XMLWrapped: &ast.Attribute{Value: token.String{Value: "parent"}},
		},
	}, false)
	a.Equal(n, "n1")
	n = parseXMLWrappedName(&ast.Param{
		Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
		Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		XML: ast.XML{
			XMLWrapped: &ast.Attribute{Value: token.String{Value: "parent"}},
		},
	}, true)
	a.Equal(n, "parent")

	// wrapped = parent>name
	n = parseXMLWrappedName(&ast.Param{
		Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
		Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		XML: ast.XML{
			XMLWrapped: &ast.Attribute{Value: token.String{Value: "parent>n"}},
		},
	}, false)
	a.Equal(n, "n")
	n = parseXMLWrappedName(&ast.Param{
		Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
		Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		XML: ast.XML{
			XMLWrapped: &ast.Attribute{Value: token.String{Value: "parent>n"}},
		},
	}, true)
	a.Equal(n, "parent")

	// wrapped = >name
	n = parseXMLWrappedName(&ast.Param{
		Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
		Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		XML: ast.XML{
			XMLWrapped: &ast.Attribute{Value: token.String{Value: ">n"}},
		},
	}, false)
	a.Equal(n, "n")
	n = parseXMLWrappedName(&ast.Param{
		Name:  &ast.Attribute{Value: token.String{Value: "n1"}},
		Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		XML: ast.XML{
			XMLWrapped: &ast.Attribute{Value: token.String{Value: ">n"}},
		},
	}, true)
	a.Empty(n)
}

func TestValidXMLParamValue(t *testing.T) {
	a := assert.New(t)

	// None
	a.NotError(validXMLValue(&ast.Param{}, "", ""))
	a.Error(validXMLValue(&ast.Param{}, "", "xx"))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNone}}}, "", ""))
	a.Error(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNone}}}, "", "xx"))

	// Number
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}}}, "", "1111"))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}}}, "", "0"))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}}}, "", "-11"))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}}}, "", "-1024.11"))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}}}, "", "-1024.1111234"))
	a.Error(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}}}, "", "fxy0"))

	// String
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}}}, "", "fxy0"))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}}}, "", ""))

	// Bool
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}}}, "", "true"))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}}}, "", "false"))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}}}, "", "1"))
	a.Error(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}}}, "", "false/true"))

	// Object
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}}}, "", ""))
	a.NotError(validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}}}, "", "{}"))

	// panic
	a.Panic(func() {
		validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "xxx"}}}, "", "{}")
	})
	a.Panic(func() {
		validXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "xxx"}}}, "", "")
	})

	// bool enum
	p := &ast.Param{
		Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}},
		Enums: []*ast.Enum{
			{Value: &ast.Attribute{Value: token.String{Value: "true"}}},
			{Value: &ast.Attribute{Value: token.String{Value: "false"}}},
		},
	}
	a.NotError(validXMLValue(p, "", "true"))

	// 不存在于枚举
	p = &ast.Param{
		Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}},
		Enums: []*ast.Enum{
			{Value: &ast.Attribute{Value: token.String{Value: "true"}}},
		},
	}
	a.Error(validXMLValue(p, "", "false"))

	// number enum
	p = &ast.Param{
		Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
		Enums: []*ast.Enum{
			{Value: &ast.Attribute{Value: token.String{Value: "1"}}},
			{Value: &ast.Attribute{Value: token.String{Value: "-1.2"}}},
		},
	}
	a.NotError(validXMLValue(p, "", "1"))
	a.NotError(validXMLValue(p, "", "-1.2"))

	// 不存在于枚举
	p = &ast.Param{
		Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
		Enums: []*ast.Enum{
			{Value: &ast.Attribute{Value: token.String{Value: "1"}}},
			{Value: &ast.Attribute{Value: token.String{Value: "-1.2"}}},
		},
	}
	a.Error(validXMLValue(p, "", "false"))
}

func TestGenXMLValue(t *testing.T) {
	a := assert.New(t)

	v := genXMLValue(testOptions, &ast.Param{})
	a.Equal(v, "")

	v = genXMLValue(testOptions, &ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNone}}})
	a.Equal(v, "")

	v = genXMLValue(testOptions, &ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}}})
	a.Equal(v, true)

	v = genXMLValue(testOptions, &ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}}})
	a.Equal(v, 1024)

	v = genXMLValue(testOptions, &ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}}})
	a.Equal(v, "1024")

	a.Panic(func() {
		v = genXMLValue(testOptions, &ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeObject}}})
	})

	a.Panic(func() {
		v = genXMLValue(testOptions, &ast.Param{Type: &ast.TypeAttribute{Value: token.String{Value: "not-exists"}}})
	})
}
