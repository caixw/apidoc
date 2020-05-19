// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"sort"

	"github.com/caixw/apidoc/v7/internal/docs/localedoc"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/node"
)

// 用于描述类型信息
type typeList struct {
	Types []*localedoc.Type
}

// NewTypes 分析 v 并将其转换成 Types 数据
func NewTypes(doc *localedoc.LocaleDoc, v interface{}) error {
	n := node.New("", reflect.ValueOf(v))
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

func (types *typeList) dumpToTypes(n *node.Node) error {
	t := &localedoc.Type{
		Name:  n.TypeName,
		Usage: localedoc.InnerXML{Text: locale.Sprintf(n.Value.Usage)},
		Items: make([]*localedoc.Item, 0, len(n.Attributes)+len(n.Elements)),
	}
	types.Types = append(types.Types, t) // 保证子元素在后显示

	for _, attr := range n.Attributes {
		appendItem(t, "@"+attr.Name, attr.Value, attr.Usage, !attr.Omitempty)

		if nn := node.New(attr.Name, attr.Value); nn.TypeName != "" && !types.typeExists(nn.TypeName) {
			if err := types.dumpToTypes(nn); err != nil {
				return err
			}
		}
	}

	for _, elem := range n.Elements {
		appendItem(t, elem.Name, elem.Value, elem.Usage, !elem.Omitempty)

		typ := node.GetRealType(elem.Type())
		v := node.GetRealValue(elem.Value)

		if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			typ = node.GetRealType(typ.Elem())
			v = reflect.New(typ).Elem()
		}

		if nn := node.New(elem.Name, v); nn.TypeName != "" && !types.typeExists(nn.TypeName) {
			if err := types.dumpToTypes(nn); err != nil {
				return err
			}
		}
	}

	if n.CData != nil {
		appendItem(t, ".", n.CData.Value, n.CData.Usage, !n.CData.Omitempty)
	}

	if n.Content != nil {
		appendItem(t, ".", n.Content.Value, n.Content.Usage, !n.Content.Omitempty)
	}

	return nil
}

func appendItem(t *localedoc.Type, name string, v reflect.Value, usageKey string, req bool) {
	var isSlice bool
	typ := node.GetRealValue(v).Type()
	for typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
		isSlice = true
		typ = typ.Elem()
	}

	tt := typ.Name()
	if vv := node.ParseValue(reflect.New(typ).Elem()); vv != nil {
		tt = vv.Name
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
