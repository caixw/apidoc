// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
)

func TestDecode(t *testing.T) {
	a := assert.New(t)

	v := &struct {
		Attr1 intTest `apidoc:"attr1,attr,usage"`
		Elem1 intTest `apidoc:"elem1,elem,usage"`
	}{}
	b := `<apidoc attr1="5"><elem1>6</elem1></apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v))
	attr1 := intTest{Value: 5, Base: Base{UsageKey: "usage", Range: core.Range{
		Start: core.Position{Character: 8},
		End:   core.Position{Character: 17},
	}}}
	elem1 := intTest{Value: 6, Base: Base{UsageKey: "usage", Range: core.Range{
		Start: core.Position{Character: 18},
		End:   core.Position{Character: 34},
	}}}
	a.Equal(v.Attr1, attr1).
		Equal(v.Elem1, elem1)

	// 数组，单个元素
	v2 := &struct {
		Attr1 intTest   `apidoc:"attr1,attr,usage"`
		Elem1 []intTest `apidoc:"elem1,elem,usage"`
	}{}
	b = `<apidoc attr1="5"><elem1>6</elem1></apidoc>`
	attr1 = intTest{Value: 5, Base: Base{UsageKey: "usage", Range: core.Range{
		Start: core.Position{Character: 8},
		End:   core.Position{Character: 17},
	}}}
	elem1 = intTest{Value: 6, Base: Base{UsageKey: "usage", Range: core.Range{
		Start: core.Position{Character: 18},
		End:   core.Position{Character: 34},
	}}}
	a.NotError(Decode(core.Block{Data: []byte(b)}, v2))
	a.Equal(v2.Attr1, attr1).
		Equal(v2.Elem1, []intTest{elem1})

	// 数组，多个元素
	v3 := &struct {
		Attr1 intTest   `apidoc:"attr1,attr,usage"`
		Elem1 []intTest `apidoc:"elem1,elem,usage"`
	}{}
	b = `<apidoc attr1="5"><elem1>6</elem1><elem1>7</elem1></apidoc>`
	attr1 = intTest{Value: 5, Base: Base{UsageKey: "usage", Range: core.Range{
		Start: core.Position{Character: 8},
		End:   core.Position{Character: 17},
	}}}
	elem1 = intTest{Value: 6, Base: Base{UsageKey: "usage", Range: core.Range{
		Start: core.Position{Character: 18},
		End:   core.Position{Character: 34},
	}}}
	elem2 := intTest{Value: 7, Base: Base{UsageKey: "usage", Range: core.Range{
		Start: core.Position{Character: 34},
		End:   core.Position{Character: 50},
	}}}
	a.NotError(Decode(core.Block{Data: []byte(b)}, v3))
	a.Equal(v3.Attr1, attr1).
		Equal(v3.Elem1, []intTest{elem1, elem2})

	// content
	v4 := &struct {
		ID      intTest `apidoc:"attr1,attr,usage"`
		Content String  `apidoc:",content"`
	}{}
	b = `<apidoc attr1="5">5555</apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v4))
	a.Equal(v4.Content, String{Value: "5555", Range: core.Range{
		Start: core.Position{Character: 18},
		End:   core.Position{Character: 22},
	}}).
		Equal(v4.ID, intTest{Value: 5, Base: Base{UsageKey: "usage", Range: core.Range{
			Start: core.Position{Character: 8},
			End:   core.Position{Character: 17},
		}}})

	// cdata
	v5 := &struct {
		Cdata *CData `apidoc:",cdata"`
	}{}
	b = `<apidoc attr1="5"><![CDATA[5555]]></apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v5))
	a.Equal(v5.Cdata, &CData{
		Value: String{Value: "5555", Range: core.Range{
			Start: core.Position{Character: 27},
			End:   core.Position{Character: 31},
		}},
		Range: core.Range{
			Start: core.Position{Character: 18},
			End:   core.Position{Character: 34},
		},
	})

	// cdata 没有围绕 CDATA，则会被忽略
	v6 := &struct {
		Cdata CData `apidoc:",cdata"`
	}{}
	b = `<apidoc attr1="5">5555</apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v6))
	a.Empty(v6.Cdata.Value.Value).True(v6.Cdata.IsEmpty())

	v7 := &struct {
		ID     *intTest    `apidoc:"id,attr,usage"`
		Name   stringTest  `apidoc:"name,elem,usage"`
		Object *objectTest `apidoc:"obj,elem,usage"`
	}{}
	b = `<apidoc id="11"><name>name</name><obj id="11"><name>n</name></obj></apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v7))
	a.Equal(v7.ID, &intTest{Value: 11, Base: Base{UsageKey: "usage", Range: core.Range{
		Start: core.Position{Character: 8},
		End:   core.Position{Character: 15},
	}}}).
		Equal(v7.Name, stringTest{Value: "name", Base: Base{UsageKey: "usage", Range: core.Range{
			Start: core.Position{Character: 16},
			End:   core.Position{Character: 33},
		}}}).
		Equal(v7.Object, &objectTest{
			ID: intTest{
				Value: 11,
				Base: Base{UsageKey: "usage", Range: core.Range{
					Start: core.Position{Character: 38},
					End:   core.Position{Character: 45},
				}},
			},
			Name: stringTest{
				Value: "n",
				Base: Base{UsageKey: "usage", Range: core.Range{
					Start: core.Position{Character: 46},
					End:   core.Position{Character: 60},
				}},
			},
			Base: Base{UsageKey: "usage", Range: core.Range{
				Start: core.Position{Character: 33},
				End:   core.Position{Character: 66},
			}},
		})

	// 多个根元素
	b = `<apidoc attr="1"></apidoc><apidoc attr="1"></apidoc>`
	a.Error(Decode(core.Block{Data: []byte(b)}, v7))

	// 多个结束元素
	b = `<apidoc attr="1"></apidoc></apidoc>`
	a.Error(Decode(core.Block{Data: []byte(b)}, v7))

	// 无效的属性值
	v8 := &struct {
		ID intTest `apidoc:"id,attr,usage"`
	}{}
	b = `<apidoc id="1xx"></apidoc></apidoc>`
	a.Error(Decode(core.Block{Data: []byte(b)}, v8))

	// StartElement.Close
	v9 := &struct {
		ID intTest `apidoc:"id,attr,usage"`
	}{}
	b = `<apidoc id="1" />`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v9))

	// 不存在的元素名
	v10 := &struct {
		ID intTest `apidoc:"id,elem,usage"`
	}{}
	b = `<apidoc id="1"><elem>11</elem></apidoc>`
	a.Panic(func() {
		Decode(core.Block{Data: []byte(b)}, v10)
	})

	// 数组元素未实现 Decoder 接口
	v11 := &struct {
		Elem []int `apidoc:"elem,elem,usage"`
	}{}
	b = `<apidoc id="1"><elem>11</elem></apidoc>`
	a.Panic(func() {
		Decode(core.Block{Data: []byte(b)}, v11)
	})
}

func TestObject_decodeAttributes(t *testing.T) {
	a := assert.New(t)

	o := &node{}
	a.NotError(o.decodeAttributes(nil))

	val := &struct {
		ID   intTest    `apidoc:"id,attr,usage"`
		Name stringTest `apidoc:"name,attr,usage"`
	}{}
	o = newNode("root", reflect.ValueOf(val))
	a.NotNil(o)
	err := o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "10"}},
		},
	})
	a.NotError(err).Equal(val.ID, intTest{Value: 10, Base: Base{UsageKey: "usage"}}).
		Equal(val.Name, stringTest{Value: "name", Base: Base{UsageKey: "usage"}})

	val = &struct {
		ID   intTest    `apidoc:"id,attr,usage"`
		Name stringTest `apidoc:"name,attr,usage"`
	}{}
	o = newNode("root", reflect.ValueOf(val))
	a.NotNil(o)
	err = o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "xx10"}},
		},
	})
	a.Error(err)

	// 带匿名成员
	val2 := &struct {
		anonymous
		ID   intTest    `apidoc:"id,attr,usage"`
		Name stringTest `apidoc:"name,attr,usage"`
	}{}
	o = newNode("root", reflect.ValueOf(val2))
	a.NotNil(o)
	err = o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "10"}},
			{Name: String{Value: "attr1"}, Value: String{Value: "11"}},
		},
	})
	a.NotError(err).
		Equal(val2.ID, intTest{Value: 10, Base: Base{UsageKey: "usage"}}).
		Equal(val2.Name, stringTest{Value: "name", Base: Base{UsageKey: "usage"}}).
		Equal(val2.Attr1, intTest{Value: 11, Base: Base{UsageKey: "usage"}})

	// 测试 AttrDecoder，返回错误
	val4 := &struct {
		ID   errIntTest `apidoc:"id,attr,usage"`
		Name stringTest `apidoc:"name,attr,usage"`
	}{}
	o = newNode("root", reflect.ValueOf(val4))
	a.NotNil(o)
	err = o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "10"}},
		},
	})
	a.Error(err)

	// 未实现 AttrDecoder
	val5 := &struct {
		ID   int        `apidoc:"id,attr,usage"`
		Name stringTest `apidoc:"name,attr,usage"`
	}{}
	o = newNode("root", reflect.ValueOf(val5))
	a.NotNil(o)
	a.Panic(func() {
		o.decodeAttributes(&StartElement{
			Attributes: []*Attribute{
				{Name: String{Value: "name"}, Value: String{Value: "name"}},
				{Name: String{Value: "id"}, Value: String{Value: "10"}},
			},
		})
	})
}

func TestFindEndElement(t *testing.T) {
	a := assert.New(t)

	p, err := NewParser(core.Block{Data: []byte("<c>1</c>")})
	a.NotError(err).NotNil(p)
	a.Error(findEndElement(p, &StartElement{Name: String{Value: "c"}}))

	p, err = NewParser(core.Block{Data: []byte("1</c>")})
	a.NotError(err).NotNil(p)
	a.NotError(findEndElement(p, &StartElement{Name: String{Value: "c"}}))

	p, err = NewParser(core.Block{Data: []byte("<c>1</c></c>")})
	a.NotError(err).NotNil(p)
	a.NotError(findEndElement(p, &StartElement{Name: String{Value: "c"}}))

	p, err = NewParser(core.Block{Data: []byte("<c attr=\">1</c></c>")})
	a.NotError(err).NotNil(p)
	a.Error(findEndElement(p, &StartElement{Name: String{Value: "c"}}))
}
