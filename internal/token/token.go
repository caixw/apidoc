// SPDX-License-Identifier: MIT

// Package token 解析 xml 内容
package token

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
)

// StartElement 表示 XML 的元素
type StartElement struct {
	core.Range
	Name       String
	Attributes []*Attribute
	Close      bool // 是否自闭合
}

// EndElement XML 的结束元素
type EndElement struct {
	core.Range
	Name String
}

// Instruction 表示 XML 的指令
type Instruction struct {
	core.Range
	Name       String
	Attributes []*Attribute
}

// Attribute 表示 XML 属性
type Attribute struct {
	core.Range
	Name  String
	Value String
}

// String 表示 XML 的字符串数据
type String struct {
	core.Range
	Value string
}

// CData 表示 XML 的 CDATA 数据
type CData struct {
	core.Range
	Value String
}

// Comment 表示 XML 的注释
type Comment struct {
	core.Range
	Value String
}

// 这些常量对应 Base 中两个相关字段的名称
const (
	usageKeyName = "UsageKey"
	rangeName    = "Range"
)

// Base 每一个 XML 节点必须包含的内容
type Base struct {
	core.Range
	UsageKey message.Reference `apidoc:"-"`
}

// Usage 返回该节点的说明内容
func (b *Base) Usage() string {
	return locale.Sprintf(b.UsageKey)
}
