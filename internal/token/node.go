// SPDX-License-Identifier: MIT

package token

import (
	"fmt"
	"reflect"
	"strings"
)

const tagName = "apidoc"

type nodeType int8

const (
	attrNode nodeType = iota
	elemNode
	cdataNode
	contentNode
)

type node struct {
	name           string // 标签名称
	attrs          []value
	elems          []value
	cdata, content value
}

type value struct {
	name string
	reflect.Value
	omitempty bool
}

func isScalar(v reflect.Value) bool {
	return v.IsValid() &&
		(v.Kind() == reflect.String || (v.Kind() >= reflect.Bool && v.Kind() <= reflect.Complex128))
}

func newNode(name string, rv reflect.Value) *node {
	rv = getRealValue(rv)
	rt := rv.Type()

	num := rv.NumField()
	if num == 0 {
		return &node{name: name}
	}

	n := &node{
		name:  name,
		attrs: make([]value, 0, num),
		elems: make([]value, 0, num),
	}

	for i := 0; i < num; i++ {
		field := rt.Field(i)
		if field.Anonymous {
			anonymous := newNode("", rv.Field(i))
			for _, attr := range anonymous.attrs {
				n.appendAttr(attr)
			}
			for _, elem := range anonymous.elems {
				n.appendElem(elem)
			}
			continue
		}
		name, node, omitempty := parseTag(field)

		v := getRealValue(rv.Field(i))
		switch node {
		case attrNode:
			n.appendAttr(value{name: name, Value: v, omitempty: omitempty})
		case elemNode:
			n.appendElem(value{name: name, Value: v, omitempty: omitempty})
		case cdataNode:
			if n.cdata.IsValid() {
				panic("已经定义了一个节点用于表示 cdata 内容")
			}
			if n.content.IsValid() {
				panic("cdata 与 content 不能同时存在")
			}
			if len(n.elems) > 0 {
				panic("cdata 与子元素不能同时存在")
			}
			n.cdata = value{name: name, Value: v, omitempty: omitempty}
		case contentNode:
			if n.content.IsValid() {
				panic("已经定义了一个节点用于表示 content 内容")
			}
			if n.cdata.IsValid() {
				panic("cdata 与 content 不能同时存在")
			}
			if len(n.elems) > 0 {
				panic("content 与子元素不能同时存在")
			}
			n.content = value{name: name, Value: v, omitempty: omitempty}
		}
	}

	return n
}

func (n *node) elem(name string) (value, bool) {
	return n.findElem(name, n.elems)
}

func (n *node) attr(name string) (value, bool) {
	return n.findElem(name, n.attrs)
}

func (n *node) findElem(name string, elems []value) (value, bool) {
	for _, e := range elems {
		if e.name == name {
			return e, true
		}
	}
	return value{}, false
}

func getRealValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		} else {
			v = v.Elem()
		}
	}

	return v
}

func (n *node) appendAttr(v value) {
	if _, found := n.attr(v.name); found {
		panic(fmt.Sprintf("存在重复的属性名称 %s", v.name))
	}
	n.attrs = append(n.attrs, v)
}

func (n *node) appendElem(v value) {
	if n.content.IsValid() || n.cdata.IsValid() {
		panic("elems 不能同时与 content 和 cdata 存在")
	}

	if _, found := n.elem(v.name); found {
		panic(fmt.Sprintf("存在重复的元素名称 %s", v.name))
	}

	n.elems = append(n.elems, v)
}

// `apidoc:"name,attr,omitempty"`
func parseTag(field reflect.StructField) (string, nodeType, bool) {
	tag := strings.TrimSpace(field.Tag.Get(tagName))
	if tag == "" {
		return field.Name, elemNode, false
	}

	props := strings.Split(tag, ",")
	switch len(props) {
	case 0:
		return field.Name, elemNode, false
	case 1:
		return strings.TrimSpace(props[0]), elemNode, false
	case 2:
		return getTagName(field, props[0]), getAttrValue(props[1]), false
	case 3:
		return getTagName(field, props[0]), getAttrValue(props[1]), getOmitemptyValue(props[2])
	default:
		panic("无效的 struct tag，最多只能有三个值")
	}
}

func getTagName(field reflect.StructField, name string) string {
	if name == "" {
		name = field.Name
	}
	return strings.TrimSpace(name)
}

func getAttrValue(v string) nodeType {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "attr":
		return attrNode
	case "elem", "":
		return elemNode
	case "cdata":
		return cdataNode
	case "content":
		return contentNode
	default:
		panic(fmt.Sprintf("无效的 struct tag: %s，第二个元素必须得是 attr、cdata 或是 elem", v))
	}
}

func getOmitemptyValue(v string) bool {
	switch strings.TrimSpace(v) {
	case "omitempty":
		return true
	case "":
		return false
	default:
		panic(fmt.Sprintf("无效的 struct tag: %s，第三个元素必须得是 omitempty 或是空值", v))
	}
}
