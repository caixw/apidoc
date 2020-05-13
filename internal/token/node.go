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
	metaNode // 用于描述节点的一些元数据
)

var stringNodeMap = map[string]nodeType{
	"attr":    attrNode,
	"elem":    elemNode,
	"cdata":   cdataNode,
	"content": contentNode,
	"meta":    metaNode,
}

// 表示一个 XML 标签节点
type node struct {
	attrs          []*value // 当前标签的属性值列表
	elems          []*value // 当前标签的元素列表
	cdata, content *value   // 当前标签如果没有子元素，则可能有普通的内容或是 CDATA 内容
	value          value    // 当前节点本身代表的值
	typeName       string   // 当前节点的类型名称
}

// 表示 XML 节点的值的反射表示方式
type value struct {
	reflect.Value
	omitempty bool
	name      string // 节点的名称

	// 当前值可能未初始化，所以保存 usage 的值，
	// 等 value 初始化之后再赋值给 Base.UsageKey
	usage string
}

func newValue(name string, v reflect.Value, omitempty bool, usage string) *value {
	return &value{
		name:      name,
		Value:     v,
		omitempty: omitempty,
		usage:     usage,
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
		return &node{value: value{name: name}}
	}

	n := &node{
		attrs: make([]*value, 0, num),
		elems: make([]*value, 0, num),
		value: value{name: name},
	}

	for i := 0; i < num; i++ {
		field := rt.Field(i)
		if field.Anonymous {
			n.appendAnonymous(rv.Field(i))
			continue
		}

		if unicode.IsLower(rune(field.Name[0])) {
			continue
		}

		fieldName, node, usage, omitempty := parseTag(field)
		if fieldName == "-" {
			continue
		}

		v := rv.Field(i)
		switch node {
		case attrNode:
			n.appendAttr(newValue(fieldName, v, omitempty, usage))
		case elemNode:
			n.appendElem(newValue(fieldName, v, omitempty, usage))
		case metaNode:
			n.typeName = fieldName
			n.value.usage = usage
			n.value.Value = rv
			if n.value.name == "" { // 顶层元素可能没有 name，此处就和 fieldName 相同
				n.value.name = fieldName
			}
		case cdataNode:
			n.setCData(newValue(fieldName, v, omitempty, usage))
		case contentNode:
			n.setContent(newValue(fieldName, v, omitempty, usage))
		}
	}

	return n
}

func (n *node) appendAnonymous(v reflect.Value) {
	anonymous := newNode("", v)

	for _, attr := range anonymous.attrs {
		n.appendAttr(attr)
	}

	for _, elem := range anonymous.elems {
		n.appendElem(elem)
	}

	if anonymous.cdata != nil {
		n.setCData(anonymous.cdata)
	}

	if anonymous.content != nil {
		n.setContent(anonymous.content)
	}
}

func (n *node) setCData(v *value) {
	if n.cdata != nil {
		panic("已经定义了一个节点用于表示 cdata 内容")
	}
	if n.content != nil {
		panic("cdata 与 content 不能同时存在")
	}
	if len(n.elems) > 0 {
		panic("cdata 与子元素不能同时存在")
	}
	n.cdata = v
}

func (n *node) setContent(v *value) {
	if n.content != nil {
		panic("已经定义了一个节点用于表示 content 内容")
	}
	if n.cdata != nil {
		panic("cdata 与 content 不能同时存在")
	}
	if len(n.elems) > 0 {
		panic("content 与子元素不能同时存在")
	}
	n.content = v
}

func (n *node) elem(name string) (*value, bool) {
	return n.findElem(name, n.elems)
}

func (n *node) attr(name string) (*value, bool) {
	return n.findElem(name, n.attrs)
}

func (n *node) findElem(name string, elems []*value) (*value, bool) {
	for _, e := range elems {
		if e.name == name {
			return e, true
		}
	}
	return nil, false
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

func (n *node) appendAttr(v *value) {
	if _, found := n.attr(v.name); found {
		panic(fmt.Sprintf("存在重复的属性名称 %s", v.name))
	}
	n.attrs = append(n.attrs, v)
}

func (n *node) appendElem(v *value) {
	if n.content != nil || n.cdata != nil {
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
		panic(fmt.Sprintf("无效的 struct tag %s:%s，数量必须介于 [3,4] 之间，当前 %d", field.Name, tag, len(props)))
	}
}

func getTagName(field reflect.StructField, name string) string {
	if name == "" {
		name = field.Name
	}
	return strings.TrimSpace(name)
}

func getNodeType(v string) nodeType {
	if node, found := stringNodeMap[strings.ToLower(strings.TrimSpace(v))]; found {
		return node
	}
	panic(fmt.Sprintf("无效的 struct tag:%s", v))
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
