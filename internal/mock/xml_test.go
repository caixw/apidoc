// SPDX-License-Identifier: MIT

package mock

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/internal/ast"
)

func TestValidXML(t *testing.T) {
	a := assert.New(t)

	for _, item := range data {
		err := validXML(item.Type, []byte(item.XML))
		a.NotError(err, "测试 %s 时返回错误 %s", item.Title, err)
	}

	p := &ast.Request{
		Name: &ast.Attribute{Value: ast.String{Value: "root"}},
		Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeObject}},
		Items: []*ast.Param{
			{
				Name: &ast.Attribute{Value: ast.String{Value: "id"}},
				Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}},
				XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
			},
			{
				Name: &ast.Attribute{Value: ast.String{Value: "desc"}},
				Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeString}},
				XML:  ast.XML{XMLExtract: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
			},
		},
	}
	content := `<root id="1024"><desc>1024</desc></root>`
	a.Error(validXML(p, []byte(content)))
}

func TestBuildXML(t *testing.T) {
	a := assert.New(t)

	for _, item := range data {
		data, err := buildXML(item.Type)
		a.NotError(err, "测试 %s 返回了错误信息 %s", item.Title, err).
			Equal(string(data), item.XML, "测试 %s 返回的数据不相等 v1:%s,v2:%s", item.Title, string(data), item.XML)
	}
}

func TestXMLValidator_find(t *testing.T) {
	item := data[len(data)-1]

	a := assert.New(t)
	v := &xmlValidator{
		param:   item.Type.Param(),
		decoder: xml.NewDecoder(strings.NewReader(item.XML)),
	}

	v.names = []string{}
	p := v.find()
	a.Nil(p)

	v.names = nil
	p = v.find()
	a.Nil(p)

	v.names = []string{}
	p = v.find()
	a.Nil(p)

	v.names = []string{"root"}
	p = v.find()
	a.NotNil(p).Equal(p.Type.V(), ast.TypeObject)

	v.names = []string{"not-exists"}
	p = v.find()
	a.Nil(p)

	v.names = []string{"root", "group", "id"}
	p = v.find()
	a.NotNil(p).Equal(p.Type.V(), ast.TypeNumber)

	v.names = []string{"root", "group", "tags", "id"}
	p = v.find()
	a.NotNil(p).Equal(p.Type.V(), ast.TypeNumber)
}

func TestValidXMLParamValue(t *testing.T) {
	a := assert.New(t)

	// None
	a.NotError(validXMLParamValue(&ast.Param{}, "", ""))
	a.Error(validXMLParamValue(&ast.Param{}, "", "xx"))
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNone}}}, "", ""))
	a.Error(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNone}}}, "", "xx"))

	// Number
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}}}, "", "1111"))
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}}}, "", "0"))
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}}}, "", "-11"))
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}}}, "", "-1024.11"))
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}}}, "", "-1024.1111234"))
	a.Error(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}}}, "", "fxy0"))

	// String
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeString}}}, "", "fxy0"))
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeString}}}, "", ""))

	// Bool
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeBool}}}, "", "true"))
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeBool}}}, "", "false"))
	a.NotError(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeBool}}}, "", "1"))
	a.Error(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeBool}}}, "", "false/true"))

	// Other
	a.Error(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeObject}}}, "", ""))
	a.Error(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeObject}}}, "", "{}"))
	a.Error(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: "xxx"}}}, "", "{}"))
	a.Error(validXMLParamValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: "xxx"}}}, "", ""))

	// bool enum
	p := &ast.Param{
		Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeBool}},
		Enums: []*ast.Enum{
			{Value: &ast.Attribute{Value: ast.String{Value: "true"}}},
			{Value: &ast.Attribute{Value: ast.String{Value: "false"}}},
		},
	}
	a.NotError(validXMLParamValue(p, "", "true"))

	// 不存在于枚举
	p = &ast.Param{
		Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeBool}},
		Enums: []*ast.Enum{
			{Value: &ast.Attribute{Value: ast.String{Value: "true"}}},
		},
	}
	a.Error(validXMLParamValue(p, "", "false"))

	// number enum
	p = &ast.Param{
		Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}},
		Enums: []*ast.Enum{
			{Value: &ast.Attribute{Value: ast.String{Value: "1"}}},
			{Value: &ast.Attribute{Value: ast.String{Value: "-1.2"}}},
		},
	}
	a.NotError(validXMLParamValue(p, "", "1"))
	a.NotError(validXMLParamValue(p, "", "-1.2"))

	// 不存在于枚举
	p = &ast.Param{
		Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}},
		Enums: []*ast.Enum{
			{Value: &ast.Attribute{Value: ast.String{Value: "1"}}},
			{Value: &ast.Attribute{Value: ast.String{Value: "-1.2"}}},
		},
	}
	a.Error(validXMLParamValue(p, "", "false"))
}

func TestGetXMLValue(t *testing.T) {
	a := assert.New(t)

	v, err := getXMLValue(&ast.Param{})
	a.NotError(err).Equal(v, "")

	v, err = getXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNone}}})
	a.NotError(err).Equal(v, "")

	v, err = getXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeBool}}})
	a.NotError(err).Equal(v, true)

	v, err = getXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeNumber}}})
	a.NotError(err).Equal(v, 1024)

	v, err = getXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeString}}})
	a.NotError(err).Equal(v, "1024")

	v, err = getXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: ast.TypeObject}}})
	a.Error(err).Nil(v)

	v, err = getXMLValue(&ast.Param{Type: &ast.TypeAttribute{Value: ast.String{Value: "not-exists"}}})
	a.Error(err).Nil(v)
}
