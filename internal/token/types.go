// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"sort"

	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/internal/docs/localedoc"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// 用于描述类型信息
type typeList struct {
	Types []*localedoc.Type
}

// NewTypes 分析 v 并将其转换成 Types 数据
func NewTypes(doc *localedoc.LocaleDoc, v interface{}, tag language.Tag) error {
	locale.SetTag(tag)

	n := newNode("", reflect.ValueOf(v))
	types := &typeList{}
	if err := types.dumpToTypes(n); err != nil {
		return err
	}

	types.sanitize()
	doc.Types = append(doc.Types, types.Types...)
	return nil
}

// 清除一些无用的数据
func (types *typeList) sanitize() {
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

func (types *typeList) dumpToTypes(n *node) error {
	t := &localedoc.Type{
		Name:  n.typeName,
		Usage: localedoc.InnerXML{Text: locale.Sprintf(n.value.usage)},
		Items: make([]*localedoc.Item, 0, len(n.attrs)+len(n.elems)),
	}
	types.Types = append(types.Types, t) // 保证子元素在后显示

	for _, attr := range n.attrs {
		appendItem(t, "@"+attr.name, attr.Value, attr.usage, !attr.omitempty)

		if nn := newNode(attr.name, attr.Value); nn.typeName != "" && !types.typeExists(nn.typeName) {
			if err := types.dumpToTypes(nn); err != nil {
				return err
			}
		}
	}

	for _, elem := range n.elems {
		appendItem(t, elem.name, elem.Value, elem.usage, !elem.omitempty)

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
		appendItem(t, ".", n.cdata.Value, n.cdata.usage, !n.cdata.omitempty)
	}

	if n.content != nil {
		appendItem(t, ".", n.content.Value, n.content.usage, !n.content.omitempty)
	}

	return nil
}

func appendItem(t *localedoc.Type, name string, v reflect.Value, usageKey string, req bool) {
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
	t.Items = append(t.Items, &localedoc.Item{
		Name:     name,
		Type:     tt,
		Required: req,
		Array:    isSlice,
		Usage:    locale.Sprintf(usageKey),
	})
}

func (types *typeList) typeExists(typeName string) bool {
	for _, t := range types.Types {
		if t.Name == typeName {
			return true
		}
	}
	return false
}
