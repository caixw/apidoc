// SPDX-License-Identifier: MIT

package node

import (
	"fmt"
	"reflect"
	"unicode"
)

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

// ParseValue 分析 v 并返回 *Value 实例
//
// 与 NewValue 的不同在于，ParseValue 会分析对象字段中是否带有 meta 的结构体标签，
// 如果有才初始化 *Value 对象，否则返回 nil。
func ParseValue(v reflect.Value) *Value {
	v = RealValue(v)
	t := v.Type()

	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("%s 的 Kind() 必须为 reflect.Struct", v.Type()))
	}

	num := t.NumField()
	for i := 0; i < num; i++ {
		field := t.Field(i)
		if field.Anonymous || unicode.IsLower(rune(field.Name[0])) || field.Tag.Get(TagName) == "-" {
			continue
		}

		if name, node, usage, omitempty := parseTag(field); node == meta {
			return NewValue(name, v, omitempty, usage)
		}
	}
	return nil
}

// IsPrimitive 是否为有效的 Go 原始类型
func IsPrimitive(v reflect.Value) bool {
	return v.IsValid() &&
		(v.Kind() == reflect.String || (v.Kind() >= reflect.Bool && v.Kind() <= reflect.Complex128))
}

// RealType 获取指针指向的类型
func RealType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// RealValue 获取指针指向的值
//
// 如果未初始化，则会对其进行初始化。
func RealValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		} else {
			v = v.Elem()
		}
	}
	return v
}
