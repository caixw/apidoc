// SPDX-License-Identifier: MIT

// Package node 处理 ast 中各个节点的结构信息
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
	Attribute Type = iota
	Element
	CData
	Content
	Meta // 用于描述节点的一些元数据
)

var stringNodeMap = map[string]Type{
	"attr":    Attribute,
	"elem":    Element,
	"cdata":   CData,
	"content": Content,
	"meta":    Meta,
}

// Node 表示一个 XML 标签节点
type Node struct {
	Attributes     []*Value // 当前标签的属性值列表
	Elements       []*Value // 当前标签的元素列表
	CData, Content *Value   // 当前标签如果没有子元素，则可能有普通的内容或是 CDATA 内容
	Value          Value    // 当前节点本身代表的值
	TypeName       string   // 当前节点的类型名称
}

// Value 表示 XML 节点的值的反射表示方式
type Value struct {
	reflect.Value
	Omitempty bool
	Name      string // 节点的名称

	// 当前值可能未初始化，所以保存 usage 的值，
	// 等 Value 初始化之后再赋值给 Base.UsageKey
	Usage string
}

// NewValue 声明 *Value 实例
func NewValue(name string, v reflect.Value, omitempty bool, usage string) *Value {
	return &Value{
		Name:      name,
		Value:     v,
		Omitempty: omitempty,
		Usage:     usage,
	}
}

// ParseValue 分析 v 中是否带有 meta 类型的字段
//
// 如果不存在会返回 nil。
func ParseValue(v reflect.Value) *Value {
	v = GetRealValue(v)
	t := v.Type()

	num := t.NumField()
	for i := 0; i < num; i++ {
		field := t.Field(i)
		if field.Anonymous || unicode.IsLower(rune(field.Name[0])) {
			continue
		}

		if name, node, usage, omitempty := parseTag(field); node == Meta {
			return NewValue(name, v, omitempty, usage)
		}
	}
	return nil
}

// IsPrimitive 是否为 go 的原始类型
func IsPrimitive(v reflect.Value) bool {
	return v.IsValid() &&
		(v.Kind() == reflect.String || (v.Kind() >= reflect.Bool && v.Kind() <= reflect.Complex128))
}

// New 声明 Node 实例
func New(name string, rv reflect.Value) *Node {
	rv = GetRealValue(rv)
	rt := GetRealType(rv.Type())

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
		case Attribute:
			n.appendAttr(NewValue(fieldName, v, omitempty, usage))
		case Element:
			n.appendElem(NewValue(fieldName, v, omitempty, usage))
		case Meta:
			n.TypeName = fieldName
			n.Value.Usage = usage
			n.Value.Value = rv
			if n.Value.Name == "" { // 顶层元素可能没有 name，此处就和 fieldName 相同
				n.Value.Name = fieldName
			}
		case CData:
			n.setCData(NewValue(fieldName, v, omitempty, usage))
		case Content:
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

// GetRealType 获取指针指向的类型
func GetRealType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// GetRealValue 获取指针指向的值
//
// 如果未初始化，则会对其进行初始化。
func GetRealValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		} else {
			v = v.Elem()
		}
	}
	return v
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
