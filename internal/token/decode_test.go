// SPDX-License-Identifier: MIT

package token

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
)

type (
	attrDecodeObject struct {
		ID int
	}
	attrDecodeInt int

	decodeObject struct {
		ID int
	}
	decodeInt int
)

var (
	_ AttrDecoder = &attrDecodeObject{}
	_ AttrDecoder = attrDecodeInt(5)

	_ Decoder = &decodeObject{}
	_ Decoder = decodeInt(5)
)

func (o *attrDecodeObject) DecodeXMLAttr(attr *Attribute) error {
	v, err := strconv.Atoi(attr.Value.Value)
	if err != nil {
		return err
	}
	o.ID = v + 1
	return nil
}

func (o attrDecodeInt) DecodeXMLAttr(attr *Attribute) error {
	_, err := strconv.Atoi(attr.Value.Value)
	if err != nil {
		return err
	}
	return nil
}

func (o *decodeObject) DecodeXML(p *Parser, start *StartElement) error {
	for {
		t, err := p.Token()
		if err != nil {
			return err
		}

		switch elem := t.(type) {
		case *String:
			v, err := strconv.Atoi(elem.Value)
			if err != nil {
				return p.WithError(elem.Start, elem.End, err)
			}
			o.ID = v
		case *EndElement:
			if elem.Name.Value != start.Name.Value {
				return p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
			}
			return nil
		default:
			panic(fmt.Sprintf("无效的节点类型 %+v", t))
		}
	}
}

func (o decodeInt) DecodeXML(p *Parser, start *StartElement) error {
	for {
		t, err := p.Token()
		if err != nil {
			return err
		}

		switch elem := t.(type) {
		case *String:
			_, err := strconv.Atoi(elem.Value)
			if err != nil {
				return p.WithError(elem.Start, elem.End, err)
			}
		case *EndElement:
			if elem.Name.Value != start.Name.Value {
				return p.NewError(elem.Start, elem.End, locale.ErrInvalidXML)
			}
			return nil
		default:
			panic(fmt.Sprintf("无效的节点类型 %+v", t))
		}
	}
}

func TestDecode(t *testing.T) {
	a := assert.New(t)

	v := &anonymous{}
	b := `<apidoc attr1="5"><elem1>6</elem1></apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v))
	a.Equal(v.Attr1, 5).
		Equal(v.Elem1, 6)

	// 数组，单个元素
	v2 := &struct {
		Attr1 int   `apidoc:"attr1,attr"`
		Elem1 []int `apidoc:"elem1,elem"`
	}{}
	b = `<apidoc attr1="5"><elem1>6</elem1></apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v2))
	a.Equal(v2.Attr1, 5).
		Equal(v2.Elem1, []int{6})

	// 数组，多个元素
	v3 := &struct {
		Attr1 int   `apidoc:"attr1,attr"`
		Elem1 []int `apidoc:"elem1,elem"`
	}{}
	b = `<apidoc attr1="5"><elem1>6</elem1><elem1>7</elem1></apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v3))
	a.Equal(v3.Attr1, 5).
		Equal(v3.Elem1, []int{6, 7})

	// content
	v4 := &struct {
		Content string `apidoc:",content"`
	}{}
	b = `<apidoc attr1="5">5555</apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v4))
	a.Equal(v4.Content, "5555")

	// cdata
	v5 := &struct {
		Cdata string `apidoc:",cdata"`
	}{}
	b = `<apidoc attr1="5"><![CDATA[5555]]></apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v5))
	a.Equal(v5.Cdata, "5555")

	// cdata 没有围绕 CDATA，则会被忽略
	v6 := &struct {
		Cdata string `apidoc:",cdata"`
	}{}
	b = `<apidoc attr1="5">5555</apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v6))
	a.Empty(v6.Cdata)

	v7 := &struct {
		ID     int        `apidoc:"id,attr"`
		Name   string     `apidoc:"name,elem"`
		Object *anonymous `apidoc:"obj,elem"`
	}{}
	b = `<apidoc id="11"><name>name</name><obj attr1="11"><elem1>11</elem1></obj></apidoc>`
	a.NotError(Decode(core.Block{Data: []byte(b)}, v7))
	a.Equal(v7.ID, 11).Equal(v7.Name, "name").
		Equal(v7.Object, &anonymous{Attr1: 11, Elem1: 11})
}

func TestObject_decodeAttributes(t *testing.T) {
	a := assert.New(t)

	o := &node{}
	a.NotError(o.decodeAttributes(nil))

	val := &struct {
		ID   int    `apidoc:"id,attr"`
		Name string `apidoc:"name,attr"`
	}{}
	o = newNode("root", reflect.ValueOf(val))
	a.NotNil(o)
	err := o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "10"}},
		},
	})
	a.NotError(err).Equal(val.ID, 10).Equal(val.Name, "name")

	// 值不能正确转换
	val = &struct {
		ID   int    `apidoc:"id,attr"`
		Name string `apidoc:"name,attr"`
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
		ID   int    `apidoc:"id,attr"`
		Name string `apidoc:"name,attr"`
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
	a.NotError(err).Equal(val2.ID, 10).Equal(val2.Name, "name").Equal(val2.Attr1, 11)

	// 测试 AttrDecoder
	val3 := &struct {
		ID   int               `apidoc:"id,attr"`
		Name string            `apidoc:"name,attr"`
		O    *attrDecodeObject `apidoc:"object,attr"`
	}{}
	o = newNode("root", reflect.ValueOf(val3))
	a.NotNil(o)
	err = o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "10"}},
			{Name: String{Value: "object"}, Value: String{Value: "11"}},
		},
	})
	a.NotError(err).Equal(val3.ID, 10).Equal(val3.Name, "name").Equal(val3.O.ID, 12)

	// 测试 AttrDecoder，返回错误
	val4 := &struct {
		ID   int               `apidoc:"id,attr"`
		Name string            `apidoc:"name,attr"`
		O    *attrDecodeObject `apidoc:"object,attr"`
	}{}
	o = newNode("root", reflect.ValueOf(val4))
	a.NotNil(o)
	err = o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "10"}},
			{Name: String{Value: "object"}, Value: String{Value: "11xx"}}, // 无法解析数值，返回错误
		},
	})
	a.Error(err)

	// 测试 AttrDecoder
	val5 := &struct {
		ID   int           `apidoc:"id,attr"`
		Name string        `apidoc:"name,attr"`
		Int  attrDecodeInt `apidoc:"int,attr"`
	}{}
	o = newNode("root", reflect.ValueOf(val5))
	a.NotNil(o)
	err = o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "10"}},
			{Name: String{Value: "int"}, Value: String{Value: "11"}},
		},
	})
	a.NotError(err).Equal(val5.ID, 10).Equal(val5.Name, "name").Equal(val5.Int, 0)

	// 测试 AttrDecoder
	val6 := &struct {
		ID   int           `apidoc:"id,attr"`
		Name string        `apidoc:"name,attr"`
		Int  attrDecodeInt `apidoc:"int,attr"`
	}{}
	o = newNode("root", reflect.ValueOf(val6))
	a.NotNil(o)
	err = o.decodeAttributes(&StartElement{
		Attributes: []*Attribute{
			{Name: String{Value: "name"}, Value: String{Value: "name"}},
			{Name: String{Value: "id"}, Value: String{Value: "10"}},
			{Name: String{Value: "int"}, Value: String{Value: "11xxx"}}, // 无法解析数值，返回错误
		},
	})
	a.Error(err)
}

func TestObject_decodeElements(t *testing.T) {
	a := assert.New(t)

	p, err := NewParser(core.Block{Data: []byte("<c>1</c>")})
	a.NotError(err).NotNil(p)
	obj := &struct {
		Content string `apidoc:"c,content"`
	}{}
	node := newNode("root", reflect.ValueOf(obj))
	a.NotError(node.decodeElements(p))
	a.Equal(obj.Content, "1")

	// 数组，多个元素
	p, err = NewParser(core.Block{Data: []byte("<c>1</c><c>2</c>")})
	a.NotError(err).NotNil(p)
	obj2 := &struct {
		Content []string `apidoc:"c,content"`
	}{}
	node = newNode("root", reflect.ValueOf(obj2))
	a.NotError(node.decodeElements(p))
	a.Equal(obj2.Content, []string{"1", "2"})

	// 数组，单个元素
	p, err = NewParser(core.Block{Data: []byte("<c>1</c>")})
	a.NotError(err).NotNil(p)
	obj2 = &struct {
		Content []string `apidoc:"c,content"`
	}{}
	node = newNode("root", reflect.ValueOf(obj2))
	a.NotError(node.decodeElements(p))
	a.Equal(obj2.Content, []string{"1"})

	// content
	p, err = NewParser(core.Block{Data: []byte("<elem>1</elem><c>1</c><c>2</c>")})
	a.NotError(err).NotNil(p)
	obj3 := &struct {
		Content []string `apidoc:"c,content"`
	}{}
	node = newNode("root", reflect.ValueOf(obj3))
	a.NotError(node.decodeElements(p))
	a.Equal(obj3.Content, []string{"1", "2"})

	// CDATA
	p, err = NewParser(core.Block{Data: []byte("<elem>1</elem><c>1</c><c>2</c><![CDATA[cdata]]>")})
	a.NotError(err).NotNil(p)
	obj4 := &struct {
		CData string `apidoc:",cdata"`
	}{}
	node = newNode("root", reflect.ValueOf(obj4))
	a.NotError(node.decodeElements(p))
	a.Equal(obj4.CData, "cdata")

	// 实现 Decoder 接口
	p, err = NewParser(core.Block{Data: []byte("<elem>1</elem><c>1</c><c>2</c><![CDATA[cdata]]>")})
	a.NotError(err).NotNil(p)
	obj5 := &struct {
		Content []*decodeObject `apidoc:"c,content"`
	}{}
	node = newNode("root", reflect.ValueOf(obj5))
	a.NotError(node.decodeElements(p))
	a.Equal(obj5.Content, []*decodeObject{{ID: 1}, {ID: 2}})

	// 实现 Decoder 接口
	p, err = NewParser(core.Block{Data: []byte("<elem>1</elem><c>1</c><![CDATA[cdata]]>")})
	a.NotError(err).NotNil(p)
	obj6 := &struct {
		Elem []decodeInt `apidoc:"elem,elem"`
	}{}
	node = newNode("root", reflect.ValueOf(obj6))
	a.NotError(node.decodeElements(p))
	a.Equal(obj6.Elem, []decodeInt{0})
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
