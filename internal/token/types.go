// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"sort"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/internal/locale"
)

// Types 所有类型信息的集合
type Types struct {
	XMLName struct{} `xml:"types"`
	Types   []*Type  `xml:"type"`
}

// Type 用于生成文档中的类型信息
type Type struct {
	XMLName struct{} `xml:"type"`

	Name  string   `xml:"name,attr"`
	Usage InnerXML `xml:"usage"`
	Items []*Item  `xml:"item,omitempty"`
}

// InnerXML 可以用于在字符串嵌套 HTML
type InnerXML struct {
	Text string `xml:",innerxml"`
}

// Item 用于描述文档类型中的单条记录内容
type Item struct {
	Name     string `xml:"name,attr"` // 变量名
	Type     string `xml:"type,attr"` // 变量的类型
	Array    bool   `xml:"array,attr"`
	Required bool   `xml:"required,attr"`
	Usage    string `xml:",innerxml"`
}

// NewTypes 分析 v，返回 Type 类型的数据
func NewTypes(v interface{}, tag language.Tag) (*Types, error) {
	locale.SetTag(tag)

	n := newNode("", reflect.ValueOf(v))
	types := &Types{}
	if err := types.dumpToTypes(n); err != nil {
		return nil, err
	}

	types.sanitize()

	return types, nil
}

// 清除一些无用的数据
func (types *Types) sanitize() {
	for _, t := range types.Types {
		if len(t.Items) == 1 && t.Items[0].Name == "." {
			t.Items = nil
		}
	}

	sort.SliceStable(types.Types, func(i, j int) bool {
		if len(types.Types[i].Items) == 0 {
			return false
		}
		return len(types.Types[j].Items) == 0
	})
}

func (types *Types) dumpToTypes(n *node) error {
	t := &Type{
		Name:  n.typeName,
		Usage: InnerXML{Text: locale.Sprintf(n.value.usage)},
		Items: make([]*Item, 0, len(n.attrs)+len(n.elems)),
	}
	types.Types = append(types.Types, t) // 保证子元素在后显示

	for _, attr := range n.attrs {
		t.appendItem("@"+attr.name, attr.Value, attr.usage, !attr.omitempty)

		if nn := newNode(attr.name, attr.Value); nn.typeName != "" && !types.typeExists(nn.typeName) {
			if err := types.dumpToTypes(nn); err != nil {
				return err
			}
		}
	}

	for _, elem := range n.elems {
		t.appendItem(elem.name, elem.Value, elem.usage, !elem.omitempty)

		typ := getRealType(elem.Type())
		v := getRealValue(elem.Value)

		if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			typ = getRealType(typ.Elem())
			v = reflect.New(typ).Elem()
		}

		if nn := newNode(elem.name, v); nn.typeName != "" && !types.typeExists(nn.typeName) {
			if err := types.dumpToTypes(nn); err != nil {
				return err
			}
		}
	}

	if n.cdata != nil {
		t.appendItem(".", n.cdata.Value, n.cdata.usage, !n.cdata.omitempty)
	}

	if n.content != nil {
		t.appendItem(".", n.content.Value, n.content.usage, !n.content.omitempty)
	}

	return nil
}

func (t *Type) appendItem(name string, v reflect.Value, usageKey string, req bool) {
	var isSlice bool
	typ := getRealValue(v).Type()
	for typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
		isSlice = true
		typ = typ.Elem()
	}

	tt := typ.Name()
	if vv := parseValue(reflect.New(typ).Elem()); vv != nil {
		tt = vv.name
	}
	t.Items = append(t.Items, &Item{
		Name:     name,
		Type:     tt,
		Required: req,
		Array:    isSlice,
		Usage:    locale.Sprintf(usageKey),
	})
}

func (types *Types) typeExists(typeName string) bool {
	for _, t := range types.Types {
		if t.Name == typeName {
			return true
		}
	}
	return false
}
