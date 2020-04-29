// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

func decode(a *assert.Assertion, xml string, v interface{}, hasErr bool) {
	p, err := NewParser(core.Block{Data: []byte(xml)})
	a.NotError(err).
		NotNil(p)

	if hasErr {
		a.Error(Decode(p, v))
		return
	}
	a.NotError(Decode(p, v))
}

func TestDecode(t *testing.T) {
	a := assert.New(t)

	v := &struct {
		Attr1 intTest `apidoc:"attr1,attr,usage"`
		Elem1 intTest `apidoc:"elem1,elem,usage"`
	}{}
	b := `<apidoc attr1="5"><elem1>6</elem1></apidoc>`
	decode(a, b, v, false)
	attr1 := intTest{Value: 5,
		Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 17},
			},
			XMLName: String{
				Value: "attr1",
				Range: core.Range{
					Start: core.Position{Character: 8},
					End:   core.Position{Character: 13},
				},
			},
		}}
	elem1 := intTest{Value: 6,
		Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 18},
				End:   core.Position{Character: 34},
			},
			XMLName: String{
				Value: "elem1",
				Range: core.Range{
					Start: core.Position{Character: 19},
					End:   core.Position{Character: 24},
				},
			},
			XMLNameEnd: String{
				Value: "elem1",
				Range: core.Range{
					Start: core.Position{Character: 28},
					End:   core.Position{Character: 33},
				},
			},
		}}
	a.Equal(v.Attr1, attr1).
		Equal(v.Elem1, elem1)

	// 自闭合标签，采用上一个相同的类型
	v = &struct {
		Attr1 intTest `apidoc:"attr1,attr,usage"`
		Elem1 intTest `apidoc:"elem1,elem,usage"`
	}{}
	b = `<apidoc attr1="5"><elem1 /></apidoc>`
	decode(a, b, v, false)
	attr1 = intTest{Value: 5,
		Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 17},
			},
			XMLName: String{
				Value: "attr1",
				Range: core.Range{
					Start: core.Position{Character: 8},
					End:   core.Position{Character: 13},
				},
			},
		}}
	elem1 = intTest{Value: 0,
		Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 18},
				End:   core.Position{Character: 27},
			},
			XMLName: String{
				Value: "elem1",
				Range: core.Range{
					Start: core.Position{Character: 19},
					End:   core.Position{Character: 24},
				},
			},
		}}
	a.Equal(v.Attr1, attr1).
		Equal(v.Elem1, elem1)

	// 数组，单个元素
	v2 := &struct {
		Attr1 intTest   `apidoc:"attr1,attr,usage"`
		Elem1 []intTest `apidoc:"elem1,elem,usage"`
	}{}
	b = `<apidoc attr1="5"><elem1>6</elem1></apidoc>`
	attr1 = intTest{Value: 5, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 8},
			End:   core.Position{Character: 17},
		},
		XMLName: String{
			Value: "attr1",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 13},
			},
		},
	}}
	elem1 = intTest{Value: 6, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 18},
			End:   core.Position{Character: 34},
		},
		XMLName: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 19},
				End:   core.Position{Character: 24},
			},
		},
		XMLNameEnd: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 28},
				End:   core.Position{Character: 33},
			},
		},
	}}
	decode(a, b, v2, false)
	a.Equal(v2.Attr1, attr1).
		Equal(v2.Elem1, []intTest{elem1})

	// 数组，多个元素
	v3 := &struct {
		Attr1 intTest   `apidoc:"attr1,attr,usage"`
		Elem1 []intTest `apidoc:"elem1,elem,usage"`
	}{}
	b = `<apidoc attr1="5"><elem1>6</elem1><elem1>7</elem1></apidoc>`
	attr1 = intTest{Value: 5, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 8},
			End:   core.Position{Character: 17},
		},
		XMLName: String{
			Value: "attr1",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 13},
			},
		},
	}}
	elem1 = intTest{Value: 6, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 18},
			End:   core.Position{Character: 34},
		},
		XMLName: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 19},
				End:   core.Position{Character: 24},
			},
		},
		XMLNameEnd: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 28},
				End:   core.Position{Character: 33},
			},
		},
	}}
	elem2 := intTest{Value: 7, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 34},
			End:   core.Position{Character: 50},
		},
		XMLName: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 35},
				End:   core.Position{Character: 40},
			},
		},
		XMLNameEnd: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 44},
				End:   core.Position{Character: 49},
			},
		},
	}}
	decode(a, b, v3, false)
	a.Equal(v3.Attr1, attr1).
		Equal(v3.Elem1, []intTest{elem1, elem2})

	// 数组，多个元素，自闭合
	v3 = &struct {
		Attr1 intTest   `apidoc:"attr1,attr,usage"`
		Elem1 []intTest `apidoc:"elem1,elem,usage"`
	}{}
	b = `<apidoc attr1="5"><elem1 /><elem1>7</elem1></apidoc>`
	attr1 = intTest{Value: 5, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 8},
			End:   core.Position{Character: 17},
		},
		XMLName: String{
			Value: "attr1",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 13},
			},
		},
	}}
	elem1 = intTest{Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 18},
			End:   core.Position{Character: 27},
		},
		XMLName: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 19},
				End:   core.Position{Character: 24},
			},
		},
	}}
	elem2 = intTest{Value: 7, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 27},
			End:   core.Position{Character: 43},
		},
		XMLName: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 28},
				End:   core.Position{Character: 33},
			},
		},
		XMLNameEnd: String{
			Value: "elem1",
			Range: core.Range{
				Start: core.Position{Character: 37},
				End:   core.Position{Character: 42},
			},
		},
	}}
	decode(a, b, v3, false)
	a.Equal(v3.Attr1, attr1).
		Equal(v3.Elem1, []intTest{elem1, elem2})

	// content
	v4 := &struct {
		ID      intTest `apidoc:"attr1,attr,usage"`
		Content String  `apidoc:",content"`
	}{}
	b = `<apidoc attr1="5">5555</apidoc>`
	decode(a, b, v4, false)
	a.Equal(v4.Content, String{Value: "5555", Range: core.Range{
		Start: core.Position{Character: 18},
		End:   core.Position{Character: 22},
	}})
	a.Equal(v4.ID, intTest{Value: 5, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 8},
			End:   core.Position{Character: 17},
		},
		XMLName: String{
			Value: "attr1",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 13},
			},
		},
	}})

	// cdata
	v5 := &struct {
		Cdata *CData `apidoc:",cdata"`
	}{}
	b = `<apidoc attr1="5"><![CDATA[5555]]></apidoc>`
	decode(a, b, v5, false)
	a.Equal(v5.Cdata, &CData{
		Value: String{Value: "5555", Range: core.Range{
			Start: core.Position{Character: 27},
			End:   core.Position{Character: 31},
		}},
		Range: core.Range{
			Start: core.Position{Character: 18},
			End:   core.Position{Character: 34},
		}},
	)

	// cdata 没有围绕 CDATA，则会被忽略
	v6 := &struct {
		Cdata CData `apidoc:",cdata"`
	}{}
	b = `<apidoc attr1="5">5555</apidoc>`
	decode(a, b, v6, false)
	a.Empty(v6.Cdata.Value.Value).True(v6.Cdata.IsEmpty())

	v7 := &struct {
		ID     *intTest    `apidoc:"id,attr,usage"`
		Name   stringTest  `apidoc:"name,elem,usage"`
		Object *objectTest `apidoc:"obj,elem,usage"`
	}{}
	b = `<apidoc id="11"><name>name</name><obj id="11"><name>n</name></obj></apidoc>`
	decode(a, b, v7, false)
	a.Equal(v7.ID, &intTest{Value: 11, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 8},
			End:   core.Position{Character: 15},
		},
		XMLName: String{
			Value: "id",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 10},
			},
		},
	}})
	a.Equal(v7.Name, stringTest{Value: "name", Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 16},
			End:   core.Position{Character: 33},
		},
		XMLName: String{
			Value: "name",
			Range: core.Range{
				Start: core.Position{Character: 17},
				End:   core.Position{Character: 21},
			},
		},
		XMLNameEnd: String{
			Value: "name",
			Range: core.Range{
				Start: core.Position{Character: 28},
				End:   core.Position{Character: 32},
			},
		},
	}})
	a.Equal(v7.Object, &objectTest{
		Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 33},
				End:   core.Position{Character: 66},
			},
			XMLName: String{
				Value: "obj",
				Range: core.Range{
					Start: core.Position{Character: 34},
					End:   core.Position{Character: 37},
				},
			},
			XMLNameEnd: String{
				Value: "obj",
				Range: core.Range{
					Start: core.Position{Character: 62},
					End:   core.Position{Character: 65},
				},
			},
		},
		ID: intTest{Value: 11, Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 38},
				End:   core.Position{Character: 45},
			},
			XMLName: String{
				Value: "id",
				Range: core.Range{
					Start: core.Position{Character: 38},
					End:   core.Position{Character: 40},
				},
			},
		}},
		Name: stringTest{Value: "n", Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 46},
				End:   core.Position{Character: 60},
			},
			XMLName: String{
				Value: "name",
				Range: core.Range{
					Start: core.Position{Character: 47},
					End:   core.Position{Character: 51},
				},
			},
			XMLNameEnd: String{
				Value: "name",
				Range: core.Range{
					Start: core.Position{Character: 55},
					End:   core.Position{Character: 59},
				},
			},
		}},
	})

	// 多个根元素
	b = `<apidoc attr="1"></apidoc><apidoc attr="1"></apidoc>`
	decode(a, b, v7, true)

	// 多个结束元素
	b = `<apidoc attr="1"></apidoc></apidoc>`
	decode(a, b, v7, true)

	// 无效的属性值
	v8 := &struct {
		ID intTest `apidoc:"id,attr,usage"`
	}{}
	b = `<apidoc id="1xx"></apidoc></apidoc>`
	decode(a, b, v8, true)

	// StartElement.Close
	v9 := &struct {
		ID intTest `apidoc:"id,attr,usage"`
	}{}
	b = `<apidoc id="1" />`
	decode(a, b, v9, false)

	// 不存在的元素名
	v10 := &struct {
		ID intTest `apidoc:"id,elem,usage"`
	}{}
	b = `<apidoc id="1"><elem>11</elem></apidoc>`
	a.Panic(func() {
		decode(a, b, v10, false)
	})

	// 数组元素未实现 Decoder 接口
	v11 := &struct {
		Elem []int `apidoc:"elem,elem,usage"`
	}{}
	b = `<apidoc id="1"><elem>11</elem></apidoc>`
	a.Panic(func() {
		decode(a, b, v11, false)
	})

	// 多个数组，未实现 Decoder 的元素
	v12 := &struct {
		Attr1 intTest       `apidoc:"attr1,attr,usage"`
		Elem1 []*objectTest `apidoc:"e,elem,usage"`
	}{}
	b = `<apidoc attr1="5">
	<e id="5"><name>6</name></e>
	<e id="7"><name>7</name></e>
</apidoc>`
	attr1 = intTest{Value: 5, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 8},
			End:   core.Position{Character: 17},
		},
		XMLName: String{
			Value: "attr1",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 13},
			},
		},
	}}
	e1 := &objectTest{
		Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 1, Line: 1},
				End:   core.Position{Character: 29, Line: 1},
			},
			XMLName: String{
				Value: "e",
				Range: core.Range{
					Start: core.Position{Character: 2, Line: 1},
					End:   core.Position{Character: 3, Line: 1},
				},
			},
			XMLNameEnd: String{
				Value: "e",
				Range: core.Range{
					Start: core.Position{Character: 27, Line: 1},
					End:   core.Position{Character: 28, Line: 1},
				},
			},
		},
		ID: intTest{
			Base: Base{
				UsageKey: "usage",
				Range: core.Range{
					Start: core.Position{Character: 4, Line: 1},
					End:   core.Position{Character: 10, Line: 1},
				},
				XMLName: String{
					Value: "id",
					Range: core.Range{
						Start: core.Position{Character: 4, Line: 1},
						End:   core.Position{Character: 6, Line: 1},
					},
				},
			},
			Value: 5,
		},
		Name: stringTest{
			Base: Base{
				UsageKey: "usage",
				Range: core.Range{
					Start: core.Position{Character: 11, Line: 1},
					End:   core.Position{Character: 25, Line: 1},
				},
				XMLName: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Character: 12, Line: 1},
						End:   core.Position{Character: 16, Line: 1},
					},
				},
				XMLNameEnd: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Character: 20, Line: 1},
						End:   core.Position{Character: 24, Line: 1},
					},
				},
			},
			Value: "6",
		},
	}
	e2 := &objectTest{
		Base: Base{
			UsageKey: "usage",
			Range: core.Range{
				Start: core.Position{Character: 1, Line: 2},
				End:   core.Position{Character: 29, Line: 2},
			},
			XMLName: String{
				Value: "e",
				Range: core.Range{
					Start: core.Position{Character: 2, Line: 2},
					End:   core.Position{Character: 3, Line: 2},
				},
			},
			XMLNameEnd: String{
				Value: "e",
				Range: core.Range{
					Start: core.Position{Character: 27, Line: 2},
					End:   core.Position{Character: 28, Line: 2},
				},
			},
		},
		ID: intTest{
			Base: Base{
				UsageKey: "usage",
				Range: core.Range{
					Start: core.Position{Character: 4, Line: 2},
					End:   core.Position{Character: 10, Line: 2},
				},
				XMLName: String{
					Value: "id",
					Range: core.Range{
						Start: core.Position{Character: 4, Line: 2},
						End:   core.Position{Character: 6, Line: 2},
					},
				},
			},
			Value: 7,
		},
		Name: stringTest{
			Base: Base{
				UsageKey: "usage",
				Range: core.Range{
					Start: core.Position{Character: 11, Line: 2},
					End:   core.Position{Character: 25, Line: 2},
				},
				XMLName: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Character: 12, Line: 2},
						End:   core.Position{Character: 16, Line: 2},
					},
				},
				XMLNameEnd: String{
					Value: "name",
					Range: core.Range{
						Start: core.Position{Character: 20, Line: 2},
						End:   core.Position{Character: 24, Line: 2},
					},
				},
			},
			Value: "7",
		},
	}
	decode(a, b, v12, false)
	a.Equal(v12.Attr1, attr1).
		Equal(2, len(v12.Elem1)).
		Equal(v12.Elem1[0], e1).
		Equal(v12.Elem1[1], e2)

	// 数组，闭合标签带属性
	type obj struct {
		Base
		ID intTest `apidoc:"id,attr,usage"`
	}
	v13 := &struct {
		Attr1 intTest `apidoc:"attr1,attr,usage"`
		Elem1 []*obj  `apidoc:"elem2,elem,usage-elem2"`
	}{}
	b = `<apidoc attr1="5"><elem2 id="6" /></apidoc>`
	attr1 = intTest{Value: 5, Base: Base{
		UsageKey: "usage",
		Range: core.Range{
			Start: core.Position{Character: 8},
			End:   core.Position{Character: 17},
		},
		XMLName: String{
			Value: "attr1",
			Range: core.Range{
				Start: core.Position{Character: 8},
				End:   core.Position{Character: 13},
			},
		},
	}}
	obj1 := &obj{
		Base: Base{
			UsageKey: "usage-elem2",
			Range: core.Range{
				Start: core.Position{Character: 18},
				End:   core.Position{Character: 34},
			},
			XMLName: String{
				Value: "elem2",
				Range: core.Range{
					Start: core.Position{Character: 19},
					End:   core.Position{Character: 24},
				},
			},
		},
		ID: intTest{
			Value: 6,
			Base: Base{
				UsageKey: "usage",
				Range: core.Range{
					Start: core.Position{Character: 25},
					End:   core.Position{Character: 31},
				},
				XMLName: String{
					Value: "id",
					Range: core.Range{
						Start: core.Position{Character: 25},
						End:   core.Position{Character: 27},
					},
				},
			},
		},
	}

	decode(a, b, v13, false)
	a.Equal(v13.Attr1, attr1)
	a.Equal(1, len(v13.Elem1)).Equal(v13.Elem1[0], obj1, "v1=%#v\nv2=%#v\n", v13.Elem1[0], obj1)

	// 闭合标签
	v14 := &struct {
		Attr1 intTest `apidoc:"attr1,attr,usage"`
		Elem1 *obj    `apidoc:"elem2,elem,usage-elem2"`
	}{}
	b = `<apidoc attr1="5"><elem2 id="60" /></apidoc>`
	decode(a, b, v14, false)
	a.NotNil(v14.Elem1).
		Equal(v14.Elem1.XMLName.Value, "elem2").
		Equal(v14.Elem1.ID.Value, 60)
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
	a.NotError(err)
	a.Equal(val.ID, intTest{Value: 10, Base: Base{
		UsageKey: "usage",
		XMLName:  String{Value: "id"},
	}})
	a.Equal(val.Name, stringTest{Value: "name", Base: Base{
		UsageKey: "usage",
		XMLName:  String{Value: "name"},
	}})

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
		Equal(val2.ID, intTest{Value: 10, Base: Base{
			UsageKey: "usage",
			XMLName:  String{Value: "id"},
		}})
	a.Equal(val2.Name, stringTest{Value: "name", Base: Base{
		UsageKey: "usage",
		XMLName:  String{Value: "name"},
	}})
	a.Equal(val2.Attr1, intTest{Value: 11, Base: Base{
		UsageKey: "usage",
		XMLName:  String{Value: "attr1"},
	}})

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
