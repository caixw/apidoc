// SPDX-License-Identifier: MIT

// Package token 解析 xml 内容
package token

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/core"
)

// 这些常量对应 Base* 中相关字段的名称
const (
	rangeName         = "Range"
	usageKeyName      = "UsageKey"
	elementTagName    = "StartTag"
	elementTagEndName = "EndTag"
	attributeNameName = "AttributeName"
)

type (
	// Name 表示 XML 中的名称
	Name struct {
		core.Range
		Prefix String
		Local  String
	}

	// StartElement 表示 XML 的元素
	StartElement struct {
		core.Range
		Name       Name
		Attributes []*Attribute
		Close      bool // 是否自闭合
	}

	// EndElement XML 的结束元素
	EndElement struct {
		core.Range
		Name Name
	}

	// Instruction 表示 XML 的指令
	Instruction struct {
		core.Range
		Name       String
		Attributes []*Attribute
	}

	// Attribute 表示 XML 属性
	Attribute struct {
		core.Range
		Name  Name
		Value String
	}

	// String 表示 XML 的字符串数据
	String struct {
		core.Range
		Value string
	}

	// CData 表示 XML 的 CDATA 数据
	CData struct {
		BaseTag
		Value String
	}

	// Comment 表示 XML 的注释
	Comment struct {
		core.Range
		Value String
	}

	// Base 所有 XML 节点的基本元素
	Base struct {
		core.Range
		UsageKey message.Reference `apidoc:"-"` // 表示对当前元素的一个说明内容的翻译 ID
	}

	// BaseAttribute 所有 XML 属性节点的基本元素
	BaseAttribute struct {
		Base
		AttributeName Name `apidoc:"-"`
	}

	// BaseTag 所有 XML 标签的基本元素
	BaseTag struct {
		Base
		StartTag Name `apidoc:"-"` // 表示起始标签名
		EndTag   Name `apidoc:"-"` // 表示标签的结束名称，如果是自闭合的标签，此值为空。
	}
)

// Match 是否与 end 相匹配
func (s *StartElement) Match(end *EndElement) bool {
	return s.Name.Equal(end.Name)
}

// Equal 两个 name 是否相等
func (n Name) Equal(v Name) bool {
	return n.Prefix.Value == v.Prefix.Value &&
		n.Local.Value == v.Local.Value
}

// String fmt.Stringer
func (n Name) String() string {
	if n.Prefix.Value == "" {
		return n.Local.Value
	}
	return n.Prefix.Value + ":" + n.Local.Value
}
