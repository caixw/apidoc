// SPDX-License-Identifier: MIT

package node

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"
)

type (
	Anonymous struct {
		Attr1 intAttr `apidoc:"attr1,attr,usage"`
		Elem1 intTag  `apidoc:"elem1,elem,usage"`
	}

	intTag struct {
		Value    int      `apidoc:"-"`
		RootName struct{} `apidoc:"number,meta,usage-number"`
	}

	intAttr struct {
		Value    int      `apidoc:"-"`
		RootName struct{} `apidoc:"number,meta,usage-number"`
	}
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
		typeName  string
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
			inputName: "attr1",
			inputNode: &struct {
				Root  struct{} `apidoc:"node,meta,usage"`
				Attr1 intAttr  `apidoc:"attr1,attr,usage"`
			}{},
			attrs:    1,
			typeName: "node",
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
				Root    struct{} `apidoc:"root,meta,usage"`
				Attr1   intAttr  `apidoc:"attr1,attr,usage"`
				Attr2   intAttr  `apidoc:"attr2,attr,usage"`
				Content string   `apidoc:",content,"`
			}{},
			attrs:    2,
			content:  true,
			typeName: "root",
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
				Cdata *string `apidoc:",cdata"`
			}{},
			attrs: 1,
			cdata: true,
		},
		{
			inputName: "elem2",
			inputNode: &struct {
				Anonymous
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
				*Anonymous
				Elem2 intTag `apidoc:"elem2,elem,usage"`
			}{
				Anonymous: &Anonymous{},
			},
			attrs: 1,
			elems: 2,
		},
		{
			inputName: "attr1_elem2_anonymous",
			inputNode: &struct {
				*Anonymous
				Elem2 *Anonymous `apidoc:"elem2,elem,usage"`
			}{
				Anonymous: &Anonymous{},
			},
			attrs: 1,
			elems: 2,
		},
	}

	for i, item := range data {
		o := New(item.inputName, reflect.ValueOf(item.inputNode))
		a.Equal(len(o.Elements), item.elems, "not equal %d\nv1=%d,v2=%d", i, len(o.Elements), item.elems).
			Equal(len(o.Attributes), item.attrs, "not equal %d\nv1=%d,v2=%d", i, len(o.Attributes), item.attrs).
			Equal(item.cdata, o.CData != nil, "not equal at %d\nv1=%v,v2=%v", i, item.cdata, o.CData != nil).
			Equal(item.content, o.Content != nil, "not equal at %d\nv1=%v,v2=%v", i, item.content, o.Content != nil).
			Equal(o.Value.Name, item.inputName)

		for k, v := range o.Attributes {
			a.True(v.IsValid(), "value.IsValid() == false 位于 %d:%d.Attributes", i, k)
		}

		for k, v := range o.Elements {
			a.True(v.IsValid(), "value.IsValid() == false 位于 %d:%d.Elements", i, k)
		}
	}

	// 数组
	o := New("empty", reflect.ValueOf(&struct {
		Elem1 []int `apidoc:"elem1,elem,usage"`
	}{}))
	a.Equal(1, len(o.Elements))

	elem, found := o.Element("elem1")
	a.True(found).True(elem.IsValid())

	elem, found = o.Element("elem1")
	a.True(found).Equal(elem.Kind(), reflect.Slice)

	// 相同的属性
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr1 int `apidoc:"attr1,attr"`
			Attr2 int `apidoc:"attr1,attr"`
			Elem1 int `apidoc:"elem1,elem"`
		}{}))
	})

	// 多个的 cdata
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr2  int     `apidoc:"attr1,attr"`
			Cdata1 *string `apidoc:",cdata"`
			Cdata2 *string `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 cdata 和 content
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr2   int     `apidoc:"attr1,attr"`
			Cdata1  *string `apidoc:",cdata"`
			Content *string `apidoc:",content"`
		}{}))
	})

	// 多个的 content
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr2 int    `apidoc:"attr1,attr"`
			C1    string `apidoc:",content"`
			C2    string `apidoc:",content"`
		}{}))
	})

	// 同时存在 cdata 和 content
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr2   int    `apidoc:"attr1,attr"`
			Content string `apidoc:",content"`
			Cdata1  string `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 cdata 和 Elements
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr2  int    `apidoc:"attr1,attr"`
			Elem1  int    `apidoc:"elem1,elem"`
			Cdata1 string `apidoc:",cdata"`
		}{}))
	})

	// 同时存在 content 和 Elements
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr2   int    `apidoc:"attr1,attr"`
			Elem1   int    `apidoc:"elem1,elem"`
			Content string `apidoc:",content"`
		}{}))
	})

	// Elements 同时与 content 和 cdata 存在
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr2   int    `apidoc:"attr1,attr"`
			Content string `apidoc:",content"`
			Elem1   int    `apidoc:"elem1,elem"`
		}{}))
	})

	// 相同的元素名
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Attr1 int `apidoc:"attr1,attr"`
			Elem1 int `apidoc:"elem1,elem"`
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})

	// 与匿名对象存在相同的元素名
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			Anonymous
			Attr1 int `apidoc:"attr1,attr"`
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})
	// 与匿名对象存在相同的元素名
	a.Panic(func() {
		New("empty", reflect.ValueOf(&struct {
			*Anonymous
			Attr1 int `apidoc:"attr1,attr"`
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})

	// 同时存在 cdata 和 Elements
	type Anonymous1 struct {
		CData string `apidoc:",cdata"`
	}
	a.Panic(func() {
		New("anonymous-content", reflect.ValueOf(&struct {
			Anonymous1
			Elem2 int `apidoc:"elem1,elem"`
		}{}))
	})

	// 同时存在两个 content
	type Anonymous2 struct {
		Content string `apidoc:",content"`
	}
	a.Panic(func() {
		New("anonymous-content", reflect.ValueOf(&struct {
			*Anonymous2
			Content string `apidoc:",content"`
		}{}))
	})
}

func TestParseTag(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		inputName string
		inputTag  string

		name      string
		node      Type
		usage     string
		omitempty bool
	}{
		{
			inputName: "Field",
			inputTag:  "field,attr,usage",

			name:  "field",
			node:  attribute,
			usage: "usage",
		},
		{
			inputName: "Field",
			inputTag:  "field,elem,usage",
			name:      "field",
			node:      element,
			usage:     "usage",
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata",
			name:      "field",
			node:      cdata,
		},
		{
			inputName: "Field",
			inputTag:  ",cdata",
			name:      "Field",
			node:      cdata,
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata,,omitempty",
			name:      "field",
			node:      cdata,
			omitempty: true,
		},
		{
			inputName: "Field",
			inputTag:  "field,cdata,,omitempty",
			name:      "field",
			node:      cdata,
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
			Tag:  reflect.StructTag(TagName + ":\"" + item.inputTag + "\""),
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
			Tag: reflect.StructTag(TagName + `:"field,not-exists"`),
		}
		parseTag(field)
	})

	// omitempty 错误
	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(TagName + `:"field,elem,usage,other"`),
		}
		parseTag(field)
	})

	// cdata 指定了 usage
	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(TagName + `:"field,cdata,usage,other"`),
		}
		parseTag(field)
	})

	// 数量太多
	a.Panic(func() {
		field := reflect.StructField{
			Tag: reflect.StructTag(TagName + `:"field,elem,usage,omitempty,xxx"`),
		}
		parseTag(field)
	})
}
