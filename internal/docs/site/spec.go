// SPDX-License-Identifier: MIT

package site

import (
	"reflect"
	"sort"

	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/node"
)

func (d *doc) newSpec(v interface{}) error {
	n := node.New("", reflect.ValueOf(v))
	if err := d.dumpToTypes(n); err != nil {
		return err
	}

	for _, t := range d.Spec {
		if len(t.Items) == 1 && t.Items[0].Name == "." {
			t.Items = nil
		}
	}

	sort.SliceStable(d.Spec, func(i, j int) bool {
		if len(d.Spec[i].Items) == 0 {
			return false
		}
		return len(d.Spec[j].Items) == 0
	})

	return nil
}

func (d *doc) dumpToTypes(n *node.Node) error {
	t := &spec{
		Name:  n.TypeName,
		Usage: innerXML{Text: locale.Sprintf(n.Value.Usage)},
		Items: make([]*item, 0, len(n.Attributes)+len(n.Elements)),
	}
	d.Spec = append(d.Spec, t) // 保证子元素在后显示

	for _, attr := range n.Attributes {
		appendItem(t, "@"+attr.Name, attr.Value, attr.Usage, !attr.Omitempty)

		if nn := node.New(attr.Name, attr.Value); nn.TypeName != "" && !d.typeExists(nn.TypeName) {
			if err := d.dumpToTypes(nn); err != nil {
				return err
			}
		}
	}

	for _, elem := range n.Elements {
		appendItem(t, elem.Name, elem.Value, elem.Usage, !elem.Omitempty)

		typ := node.RealType(elem.Type())
		v := node.RealValue(elem.Value)

		if typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
			typ = node.RealType(typ.Elem())
			v = reflect.New(typ).Elem()
		}

		if nn := node.New(elem.Name, v); nn.TypeName != "" && !d.typeExists(nn.TypeName) {
			if err := d.dumpToTypes(nn); err != nil {
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

func appendItem(t *spec, name string, v reflect.Value, usageKey string, req bool) {
	var isSlice bool
	typ := node.RealValue(v).Type()
	for typ.Kind() == reflect.Slice || typ.Kind() == reflect.Array {
		isSlice = true
		typ = typ.Elem()
	}

	tt := typ.Name()
	if vv := node.ParseValue(reflect.New(typ).Elem()); vv != nil {
		tt = vv.Name
	}
	t.Items = append(t.Items, &item{
		Name:     name,
		Type:     tt,
		Required: req,
		Array:    isSlice,
		Usage:    locale.Sprintf(usageKey),
	})
}

func (d *doc) typeExists(typeName string) bool {
	for _, t := range d.Spec {
		if t.Name == typeName {
			return true
		}
	}
	return false
}
