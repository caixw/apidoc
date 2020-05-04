// SPDX-License-Identifier: MIT

package token

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

const tagName = "apidoc"

type nodeType int8

const (
	attrNode nodeType = iota
	elemNode
	cdataNode
	contentNode
)

var (
	cdataType   = reflect.TypeOf(CData{})
	contentType = reflect.TypeOf(String{})
)

// 表示一个 XML 标签节点
type node struct {
	name           string  // 标签名称
	attrs          []value // 当前标签的属性值列表
	elems          []value // 当前标签的元素列表
	cdata, content value   // 当前标签如果没有子元素，则可能有普通的内容或是 CDATA 内容
}

// 表示 XML 节点的值的反射表示方式
type value struct {
	reflect.Value
	zero      interface{} // 当前类型的零值
	name      string      // 节点的名称
	omitempty bool

	// 当前值可能未初始化，所以保存 usage 的值，
	// 等 value 初始化之后再赋值给 Base.UsageKey
	usage string
}

func parseRootElement(v interface{}) value {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		panic("参数 v 不是一个有效的根元素")
	}

	root, found := getRealType(rv.Type()).FieldByName(rootElementTagName)
	if !found {
		panic(fmt.Sprintf("根元素 %s 未指定 %s 字段", getRealType(rv.Type()), rootElementTagName))
	}

	name, _, usage, omitempty := parseTag(root)
	if name == "-" {
		panic(fmt.Sprintf("根元素 %s.%s 的标签值 %s 不能为 -", rv.Type(), rootElementTagName, tagName))
	}

	return initValue(name, rv, omitempty, usage)
}

func initValue(name string, v reflect.Value, omitempty bool, usage string) value {
	return value{
		name:      name,
		Value:     v,
		zero:      reflect.Zero(v.Type()).Interface(),
		omitempty: omitempty,
		usage:     usage,
	}
}

func (v value) isOmitempty() bool {
	if !v.omitempty {
		return false
	}

	switch v.Value.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	default:
		return v.zero == v.Interface()
	}
}

func isPrimitive(v reflect.Value) bool {
	return v.IsValid() &&
		(v.Kind() == reflect.String || (v.Kind() >= reflect.Bool && v.Kind() <= reflect.Complex128))
}

func newNode(name string, rv reflect.Value) *node {
	rv = getRealValue(rv)
	rt := rv.Type()

	num := rt.NumField()
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

		if unicode.IsLower(rune(field.Name[0])) || field.Name == rootElementTagName {
			continue
		}

		name, node, usage, omitempty := parseTag(field)
		if name == "-" {
			continue
		}

		v := rv.Field(i)
		switch node {
		case attrNode:
			n.appendAttr(initValue(name, v, omitempty, usage))
		case elemNode:
			n.appendElem(initValue(name, v, omitempty, usage))
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
			if getRealType(field.Type) != cdataType {
				panic("cdata 的类型只能是 *CData")
			}
			n.cdata = initValue(name, v, omitempty, usage)
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
			if getRealType(field.Type) != contentType {
				panic("content 的类型只能是 *String")
			}
			n.content = initValue(name, v, omitempty, usage)
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

func getRealType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
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

// `apidoc:"name,attr,usage,omitempty"`
func parseTag(field reflect.StructField) (string, nodeType, string, bool) {
	tag := strings.TrimSpace(field.Tag.Get(tagName))
	if tag == "-" {
		return "-", 0, "", false
	}

	props := strings.Split(tag, ",")
	switch len(props) {
	case 2:
		return getTagName(field, props[0]), getNodeType(props[1]), "", false
	case 3:
		node := getNodeType(props[1])
		return getTagName(field, props[0]), node, getUsageKey(node, props[2]), false
	case 4:
		node := getNodeType(props[1])
		return getTagName(field, props[0]), node, getUsageKey(node, props[2]), getOmitempty(props[3])
	default:
		panic(fmt.Sprintf("无效的 struct tag %s，数量必须介于 [3,4] 之间", field.Name))
	}
}

func getTagName(field reflect.StructField, name string) string {
	if name == "" {
		name = field.Name
	}
	return strings.TrimSpace(name)
}

func getNodeType(v string) nodeType {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "attr":
		return attrNode
	case "elem":
		return elemNode
	case "cdata":
		return cdataNode
	case "content":
		return contentNode
	default:
		panic("无效的 struct tag，第二个元素必须得是 attr、cdata、content 或是 elem")
	}
}

func getOmitempty(v string) bool {
	switch strings.TrimSpace(v) {
	case "omitempty":
		return true
	case "":
		return false
	default:
		panic("无效的 struct tag，第四个元素必须得是 omitempty 或是空值")
	}
}

func getUsageKey(node nodeType, v string) string {
	need := (node != cdataNode) && (node != contentNode)

	if v == "" && need {
		panic("无效的 struct tag，当类型为 cdata 和 content，不能指定 usage 属性")
	} else if v != "" && !need {
		panic("无效的 struct tag，当类型不为 cdata 和 content，必须指定 usage 属性")
	}
	return v
}
