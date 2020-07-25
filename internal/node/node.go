// SPDX-License-Identifier: MIT

// Package node 处理 ast 中各个节点的结构信息
//
// struct tag
//
// 标签属性分为 4 个字段，其中前三个是必填的：
//  apidoc:"name,node-type,usage-key,omitempty"
// name 表示当前标签的名称，或是节点表示的类型；
// node-type 表示当前节点的类型，可以是以下值：
//  - elem 表示这是一个子元素；
//  - attr 表示为一个 XML 属性；
//  - cdata 表示为 CDATA 数据；
//  - content 表示为普通的字符串值；
//  - meta 表示这个字段仅用于描述当前元素的元数据，比如元素的名称等；
// usage-key 指定了当前元素的翻译项；
// omitempty 表示当前值为空时，是否可以忽略。
package node

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// TagName 结构体标签的名称
const TagName = "apidoc"

// Type 节点类型
type Type int8

// 节点类型的值
const (
	attribute Type = iota
	element
	cdata
	content
	meta // 用于描述节点的一些元数据
)

var stringNodeMap = map[string]Type{
	"attr":    attribute,
	"elem":    element,
	"cdata":   cdata,
	"content": content,
	"meta":    meta,
}

// Node 表示一个 XML 标签节点
type Node struct {
	Attributes     []*Value // 当前标签的属性值列表
	Elements       []*Value // 当前标签的元素列表
	CData, Content *Value   // 当前标签如果没有子元素，则可能有普通的内容或是 CDATA 内容
	Value          Value    // 当前节点本身代表的值
	TypeName       string   // 当前节点的类型名称
}

// New 声明 Node 实例
func New(name string, rv reflect.Value) *Node {
	rv = RealValue(rv)
	rt := rv.Type()

	num := rt.NumField()
	if num == 0 {
		return &Node{Value: Value{Name: name}}
	}

	n := &Node{
		Attributes: make([]*Value, 0, num),
		Elements:   make([]*Value, 0, num),
		Value:      Value{Name: name},
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
		case attribute:
			n.appendAttr(NewValue(fieldName, v, omitempty, usage))
		case element:
			n.appendElem(NewValue(fieldName, v, omitempty, usage))
		case meta:
			n.TypeName = fieldName
			n.Value.Usage = usage
			n.Value.Value = rv
			if n.Value.Name == "" { // 顶层元素可能没有 name，此处就和 fieldName 相同
				n.Value.Name = fieldName
			}
		case cdata:
			n.setCData(NewValue(fieldName, v, omitempty, usage))
		case content:
			n.setContent(NewValue(fieldName, v, omitempty, usage))
		}
	}

	return n
}

func (n *Node) appendAnonymous(v reflect.Value) {
	anonymous := New("", v)

	for _, attr := range anonymous.Attributes {
		n.appendAttr(attr)
	}

	for _, elem := range anonymous.Elements {
		n.appendElem(elem)
	}

	if anonymous.CData != nil {
		n.setCData(anonymous.CData)
	}

	if anonymous.Content != nil {
		n.setContent(anonymous.Content)
	}
}

func (n *Node) setCData(v *Value) {
	if n.CData != nil {
		panic("已经定义了一个节点用于表示 cdata 内容")
	}
	if n.Content != nil {
		panic("cdata 与 content 不能同时存在")
	}
	if len(n.Elements) > 0 {
		panic("cdata 与子元素不能同时存在")
	}
	n.CData = v
}

func (n *Node) setContent(v *Value) {
	if n.Content != nil {
		panic("已经定义了一个节点用于表示 content 内容")
	}
	if n.CData != nil {
		panic("cdata 与 content 不能同时存在")
	}
	if len(n.Elements) > 0 {
		panic("content 与子元素不能同时存在")
	}
	n.Content = v
}

// Element 查找名称为 name 的节点元素
func (n *Node) Element(name string) (*Value, bool) {
	return n.findElem(name, n.Elements)
}

// Attribute 查找名称为 name 的节点属性
func (n *Node) Attribute(name string) (*Value, bool) {
	return n.findElem(name, n.Attributes)
}

func (n *Node) findElem(name string, elems []*Value) (*Value, bool) {
	for _, e := range elems {
		if e.Name == name {
			return e, true
		}
	}
	return nil, false
}

func (n *Node) appendAttr(v *Value) {
	if _, found := n.Attribute(v.Name); found {
		panic(fmt.Sprintf("存在重复的属性名称 %s", v.Name))
	}
	n.Attributes = append(n.Attributes, v)
}

func (n *Node) appendElem(v *Value) {
	if n.Content != nil || n.CData != nil {
		panic("Elements 不能同时与 content 和 cdata 存在")
	}

	if _, found := n.Element(v.Name); found {
		panic(fmt.Sprintf("存在重复的元素名称 %s", v.Name))
	}

	n.Elements = append(n.Elements, v)
}

// `apidoc:"name,attr,usage,omitempty"`
func parseTag(field reflect.StructField) (string, Type, string, bool) {
	tag := strings.TrimSpace(field.Tag.Get(TagName))
	if tag == "-" {
		return "-", 0, "", false
	}

	props := strings.Split(tag, ",")
	switch len(props) {
	case 2:
		return getTagName(field, props[0]), getNodeType(props[1]), "", false
	case 3:
		node := getNodeType(props[1])
		return getTagName(field, props[0]), node, props[2], false
	case 4:
		node := getNodeType(props[1])
		return getTagName(field, props[0]), node, props[2], getOmitempty(props[3])
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

func getNodeType(v string) Type {
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
