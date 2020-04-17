// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"
)

type anonymous struct {
	Attr1 int `apidoc:"attr1,attr"`
	Elem1 int `apidoc:"elem1,elem"`
}

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
				Attr1 int `apidoc:"attr1,attr"`
			}{},
			attrs: 1,
		},
		{
			inputName: "attr2",
			inputNode: &struct {
				Attr1 int `apidoc:"attr1,attr"`
				Attr2 int `apidoc:"attr2,attr"`
			}{},
			attrs: 2,
		},
		{
			inputName: "attr1_attr2_array",
			inputNode: &struct {
				Attr1 int   `apidoc:"attr1,attr"`
				Attr2 []int `apidoc:"attr2,attr"`
			}{},
			attrs: 2,
		},
		{
			inputName: "attr2_content",
			inputNode: &struct {
				Attr1   int `apidoc:"attr1,attr"`
				Attr2   int `apidoc:"attr2,attr"`
				Content int `apidoc:",content"`
			}{},
			attrs:   2,
			content: true,
		},
		{
			inputName: "attr1_elem1",
			inputNode: &struct {
				Attr1 int `apidoc:"attr1,attr"`
				Elem1 int `apidoc:"elem1,elem"`
			}{},
			attrs: 1,
			elems: 1,
		},
		{
			inputName: "attr1_elem2",
			inputNode: &struct {
				Attr1 int `apidoc:"attr1,attr"`
				Elem1 int `apidoc:"elem1,elem"`
				Elem2 int `apidoc:"elem2,elem"`
			}{},
			attrs: 1,
			elems: 2,
		},
		{
			inputName: "attr1_elem2_cdata",
			inputNode: &struct {
				Attr1 int `apidoc:"attr1,attr"`
				Cdata int `apidoc:",cdata"`
			}{},
			attrs: 1,
			cdata: true,
		},
		{
			inputName: "attr1_elem2_anonymous",
			inputNode: &struct {
				anonymous
				Elem2 int `apidoc:"elem2,elem"`
			}{},
			attrs: 1,
			elems: 2,
		},
		{
			inputName: "attr1_elem2_anonymous",
			inputNode: &struct {
				*anonymous
				Elem2 int `apidoc:"elem2,elem"`
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
				Elem2 *anonymous `apidoc:"elem2,elem"`
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
			Equal(o.name, item.inputName)

		for k, v := range o.attrs {
			a.True(v.IsValid(), "value.IsValid() == false 位于 %d:%s.attrs", i, k)
			a.True(v.Kind() != reflect.Ptr, "非法的指针类型位于 %d:%s.attrs", i, k)
		}

		for k, v := range o.elems {
			a.True(v.IsValid(), "value.IsValid() == false 位于 %d:%s.elems", i, k)
			a.True(v.Kind() != reflect.Ptr, "非法的指针类型位于 %d:%s.elems", i, k)
		}
	}

	// 数组
	o := newNode("empty", reflect.ValueOf(&struct {
		Elem1 []int `apidoc:"elem1,elem"`
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

	// 相同的 cdata
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2  int `apidoc:"attr1,attr"`
			Cdata1 int `apidoc:",cdata"`
			Cdata2 int `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 cdata 和 content
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2   int `apidoc:"attr1,attr"`
			Cdata1  int `apidoc:",cdata"`
			Content int `apidoc:",content"`
		}{}))
	})

	// 相同的 content
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2 int `apidoc:"attr1,attr"`
			C1    int `apidoc:",content"`
			C2    int `apidoc:",content"`
		}{}))
	})

	// 同时存在 cdata 和 content
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2   int `apidoc:"attr1,attr"`
			Content int `apidoc:",content"`
			Cdata1  int `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 cdata 和 elems
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2  int `apidoc:"attr1,attr"`
			Elem1  int `apidoc:"elem1,elem"`
			Cdata1 int `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 content 和 elems
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2   int `apidoc:"attr1,attr"`
			Elem1   int `apidoc:"elem1,elem"`
			Content int `apidoc:",content"`
		}{}))
	})

	// elems 同时与 content 和 cdata 存在
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			Attr2   int `apidoc:"attr1,attr"`
			Content int `apidoc:",content"`
			Elem1   int `apidoc:"elem1,elem"`
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

	// 相同的元素名
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			anonymous
			Attr1 int `apidoc:"attr1,attr"`
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})
	// 相同的元素名
	a.Panic(func() {
		newNode("empty", reflect.ValueOf(&struct {
			*anonymous
			Attr1 int `apidoc:"attr1,attr"`
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})
}

func TestParseTag(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		inputName string
		inputTag  string
		name      string
		node      nodeType
		omitempty bool
	}{
		{
			inputName: "Field",
			name:      "Field",
			node:      elemNode,
		},
		{ // 全部采用默认属性
			inputName: "Field",
			inputTag:  ",,",
			name:      "Field",
			node:      elemNode,
		},
		{
			inputName: "Field",
			inputTag:  "field",
			name:      "field",
			node:      elemNode,
		},
		{
			inputName: "Field",
			inputTag:  "field,attr",
			name:      "field",
			node:      attrNode,
		},
		{
			inputName: "Field",
			inputTag:  "field,elem",
			name:      "field",
			node:      elemNode,
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata",
			name:      "field",
			node:      cdataNode,
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata,",
			name:      "field",
			node:      cdataNode,
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata,omitempty",
			name:      "field",
			node:      cdataNode,
			omitempty: true,
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata,omitempty",
			name:      "field",
			node:      cdataNode,
			omitempty: true,
		},
		{
			inputName: "Field",
			inputTag:  ",cdata",
			name:      "Field",
			node:      cdataNode,
		},
		{
			inputName: "Field",
			inputTag:  ",",
			name:      "Field",
			node:      elemNode,
		},
		{
			inputName: "Field",
			inputTag:  ",content",
			name:      "Field",
			node:      contentNode,
		},
	}

	for i, item := range data {
		field := reflect.StructField{
			Name: item.inputName,
			Tag:  reflect.StructTag(tagName + ":\"" + item.inputTag + "\""),
		}

		name, node, omitempty := parseTag(field)
		a.Equal(name, item.name, "not equal at %d\nv1=%+v,v2=%+v", i, name, item.name).
			Equal(node, item.node, "not equal at %d\nv1=%+v\nv2=%+v", i, node, item.node).
			Equal(omitempty, item.omitempty, "not equal at %d\nv1=%+v\nv2=%+v", i, omitempty, item.omitempty)
	}

	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(tagName + `:"field,not-exists"`),
		}
		parseTag(field)
	})

	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(tagName + `:"field,elem,other"`),
		}
		parseTag(field)
	})

	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(tagName + `:"field,elem,omitempty,xxx"`),
		}
		parseTag(field)
	})
}
