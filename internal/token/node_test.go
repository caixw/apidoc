// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"
)

func TestNewNode(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		inputName string
		inputNode interface{}
		elems     int  // 子元素的数量
		attrs     int  // 属性的数量
		cdata     bool // cdata 的标签名称
		content   bool // content 的标签名称
	}{
		{
			inputName: "empty",
			inputNode: &struct{}{},
		},
		{
			inputName: "attr1",
			inputNode: &struct {
				Attr1 intAttr `apidoc:"attr1,attr,usage"`
			}{},
			attrs: 1,
		},
		{
			inputName: "attr2",
			inputNode: &struct {
				Attr1 intAttr `apidoc:"attr1,attr,usage"`
				Attr2 intAttr `apidoc:"attr2,attr,usage"`
			}{},
			attrs: 2,
		},
		{
			inputName: "attr1",
			inputNode: &struct {
				Attr1 intTag  `apidoc:"-"`
				Attr2 intAttr `apidoc:"attr2,attr,usage"`
			}{},
			attrs: 1,
		},
		{
			inputName: "attr1_attr2_array",
			inputNode: &struct {
				Attr1 intAttr   `apidoc:"attr1,attr,usage"`
				Attr2 []intAttr `apidoc:"attr2,attr,usage"`
			}{},
			attrs: 2,
		},
		{
			inputName: "attr2_content",
			inputNode: &struct {
				Attr1   intAttr `apidoc:"attr1,attr,usage"`
				Attr2   intAttr `apidoc:"attr2,attr,usage"`
				Content String  `apidoc:",content,"`
			}{},
			attrs:   2,
			content: true,
		},
		{
			inputName: "attr1_elem1",
			inputNode: &struct {
				Attr1 intAttr `apidoc:"attr1,attr,usage"`
				Elem1 intTag  `apidoc:"elem1,elem,usage"`
			}{},
			attrs: 1,
			elems: 1,
		},
		{
			inputName: "attr1_elem2",
			inputNode: &struct {
				Attr1 intAttr `apidoc:"attr1,attr,usage"`
				Elem1 intTag  `apidoc:"elem1,elem,usage"`
				Elem2 intTag  `apidoc:"elem2,elem,usage"`
			}{},
			attrs: 1,
			elems: 2,
		},
		{
			inputName: "attr1_elem1",
			inputNode: &struct {
				Attr1 intAttr `apidoc:"attr1,attr,usage"`
				Elem1 intTag  `apidoc:"-"`
				Elem2 intTag  `apidoc:"elem2,elem,usage"`
			}{},
			attrs: 1,
			elems: 1,
		},
		{
			inputName: "attr1_cdata",
			inputNode: &struct {
				Attr1 intAttr `apidoc:"attr1,attr,usage"`
				Cdata *CData  `apidoc:",cdata"`
			}{},
			attrs: 1,
			cdata: true,
		},
		{
			inputName: "elem2",
			inputNode: &struct {
				anonymous
				Elem2 intTag `apidoc:"elem2,elem,usage"`
			}{},
			attrs: 1,
			elems: 2,
		},
		{ // 包含小字名称的字段
			inputName: "elem1",
			inputNode: &struct {
				elem1 intTag `apidoc:"elem1,elem,usage"`
				Elem2 intTag `apidoc:"elem2,elem,usage"`
			}{},
			elems: 1,
		},
		{ // 匿名字段小写不受影响
			inputName: "attr1_elem2_anonymous",
			inputNode: &struct {
				*anonymous
				Elem2 intTag `apidoc:"elem2,elem,usage"`
			}{
				anonymous: &anonymous{},
			},
			attrs: 1,
			elems: 2,
		},
		{
			inputName: "attr1_elem2_anonymous",
			inputNode: &struct {
				*anonymous
				Elem2 *anonymous `apidoc:"elem2,elem,usage"`
			}{
				anonymous: &anonymous{},
			},
			attrs: 1,
			elems: 2,
		},
	}

	for i, item := range data {
		o := newNode(item.inputName, reflect.ValueOf(item.inputNode))
		a.Equal(len(o.elems), item.elems, "not equal %d\nv1=%d,v2=%d", i, len(o.elems), item.elems).
			Equal(len(o.attrs), item.attrs, "not equal %d\nv1=%d,v2=%d", i, len(o.attrs), item.attrs).
			Equal(item.cdata, o.cdata.IsValid(), "not equal at %d\nv1=%v,v2=%v", i, item.cdata, o.cdata.IsValid()).
			Equal(item.content, o.content.IsValid(), "not equal at %d\nv1=%v,v2=%v", i, item.content, o.content.IsValid()).
			Equal(o.value.name, item.inputName)

		for k, v := range o.attrs {
			a.True(v.IsValid(), "value.IsValid() == false 位于 %d:%d.attrs", i, k)
		}

		for k, v := range o.elems {
			a.True(v.IsValid(), "value.IsValid() == false 位于 %d:%d.elems", i, k)
		}
	}

	// 数组
	o := newNode("empty", reflect.ValueOf(&struct {
		Elem1 []int `apidoc:"elem1,elem,usage"`
	}{}))
	a.Equal(1, len(o.elems))

	elem, found := o.elem("elem1")
	a.True(found).True(elem.IsValid())

	elem, found = o.elem("elem1")
	a.True(found).Equal(elem.Kind(), reflect.Slice)

	// 相同的属性
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr1 int `apidoc:"attr1,attr"`
			Attr2 int `apidoc:"attr1,attr"`
			Elem1 int `apidoc:"elem1,elem"`
		}{}))
	})

	// 多个的 cdata
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2  int    `apidoc:"attr1,attr"`
			Cdata1 *CData `apidoc:",cdata"`
			Cdata2 *CData `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 cdata 和 content
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2   int     `apidoc:"attr1,attr"`
			Cdata1  *CData  `apidoc:",cdata"`
			Content *String `apidoc:",content"`
		}{}))
	})

	// 多个的 content
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2 int    `apidoc:"attr1,attr"`
			C1    String `apidoc:",content"`
			C2    String `apidoc:",content"`
		}{}))
	})

	// 同时存在 cdata 和 content
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2   int    `apidoc:"attr1,attr"`
			Content String `apidoc:",content"`
			Cdata1  CData  `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 cdata 和 elems
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2  int   `apidoc:"attr1,attr"`
			Elem1  int   `apidoc:"elem1,elem"`
			Cdata1 CData `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 content 和 elems
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2   int    `apidoc:"attr1,attr"`
			Elem1   int    `apidoc:"elem1,elem"`
			Content String `apidoc:",content"`
		}{}))
	})

	// elems 同时与 content 和 cdata 存在
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2   int    `apidoc:"attr1,attr"`
			Content String `apidoc:",content"`
			Elem1   int    `apidoc:"elem1,elem"`
		}{}))
	})

	// 相同的元素名
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr1 int `apidoc:"attr1,attr"`
			Elem1 int `apidoc:"elem1,elem"`
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})

	// 与匿名对象存在相同的元素名
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			anonymous
			Attr1 int `apidoc:"attr1,attr"`
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})
	// 与匿名对象存在相同的元素名
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			*anonymous
			Attr1 int `apidoc:"attr1,attr"`
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})
}

func TestNode_isOmitempty(t *testing.T) {
	a := assert.New(t)

	v := value{omitempty: false}
	a.False(v.isOmitempty())

	v = initValue("elem", reflect.ValueOf(int(0)), true, "usage")
	a.True(v.isOmitempty())
	v.Value = reflect.ValueOf(int(5))
	a.False(v.isOmitempty())

	v.Value = reflect.ValueOf(uint(0))
	a.True(v.isOmitempty())
	v.Value = reflect.ValueOf(uint(5))
	a.False(v.isOmitempty())

	v.Value = reflect.ValueOf(float64(0))
	a.True(v.isOmitempty())
	v.Value = reflect.ValueOf(float32(5))
	a.False(v.isOmitempty())

	v.Value = reflect.ValueOf([]byte{})
	a.True(v.isOmitempty())
	v.Value = reflect.ValueOf([]byte{0})
	a.False(v.isOmitempty())

	v.Value = reflect.ValueOf(false)
	a.True(v.isOmitempty())
	v.Value = reflect.ValueOf(true)
	a.False(v.isOmitempty())

	v.Value = reflect.ValueOf(map[string]string{})
	a.True(v.isOmitempty())
	v.Value = reflect.ValueOf(map[string]string{"id": "0"})
	a.False(v.isOmitempty())
}

func TestParseTag(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		inputName string
		inputTag  string

		name      string
		node      nodeType
		usage     string
		omitempty bool
	}{
		{
			inputName: "Field",
			inputTag:  "field,attr,usage",

			name:  "field",
			node:  attrNode,
			usage: "usage",
		},
		{
			inputName: "Field",
			inputTag:  "field,elem,usage",
			name:      "field",
			node:      elemNode,
			usage:     "usage",
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata",
			name:      "field",
			node:      cdataNode,
		},
		{
			inputName: "Field",
			inputTag:  ",cdata",
			name:      "Field",
			node:      cdataNode,
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata,,omitempty",
			name:      "field",
			node:      cdataNode,
			omitempty: true,
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata,,omitempty",
			name:      "field",
			node:      cdataNode,
			omitempty: true,
		},
		{
			inputName: "Field",
			inputTag:  "-",
			name:      "-",
			node:      0,
			usage:     "",
		},
	}

	for i, item := range data {
		field := reflect.StructField{
			Name: item.inputName,
			Tag:  reflect.StructTag(tagName + ":\"" + item.inputTag + "\""),
		}

		name, node, usage, omitempty := parseTag(field)
		a.Equal(name, item.name, "not equal at %d\nv1=%+v,v2=%+v", i, name, item.name).
			Equal(node, item.node, "not equal at %d\nv1=%+v\nv2=%+v", i, node, item.node).
			Equal(usage, item.usage, "not equal at %d\nv1=%+v\nv2=%+v", i, usage, item.usage).
			Equal(omitempty, item.omitempty, "not equal at %d\nv1=%+v\nv2=%+v", i, omitempty, item.omitempty)
	}

	// 数量不够
	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(tagName + `:"field,not-exists"`),
		}
		parseTag(field)
	})

	// omitempty 错误
	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(tagName + `:"field,elem,usage,other"`),
		}
		parseTag(field)
	})

	// cdata 指定了 usage
	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(tagName + `:"field,cdata,usage,other"`),
		}
		parseTag(field)
	})

	// elem 未指定 usage
	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(tagName + `:"field,elem,"`),
		}
		parseTag(field)
	})

	// 数量太多
	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(tagName + `:"field,elem,usage,omitempty,xxx"`),
		}
		parseTag(field)
	})
}
