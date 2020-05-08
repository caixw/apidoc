// SPDX-License-Identifier: MIT

// Package token 解析 xml 内容
//
// struct tag
//
// 标签属性分为 4 个字段，其中前三个是必填的：
//  apidoc:"name,node-type,usage-key,omitempty"
// 第一个标签指定名称，如果为空，则直接采用字段的名称；
// 第二个标签指定标签的类型，可以是 elem 表示子元素，attr 表示属性，
// cdata 表示当前 XML 元素子元素内容 CDATA 数据，content
// 表示当前 XML 的子元素内容作为字符串保存至 content；
// 第三个元素用于指定当前元素的使用说明的本地化 ID，
// 加载后调用相关的方法会被翻译成本地化的语言内容返回；
// 第四个参数表示是否忽略空值，与标准库的 omitempty 相同功能，默认为 false。
//
// 根对象
//
// 根对象必须添加一个 RootName 字段指定根名称以其它属性：
//  type Root struct {
//      RootName struct{} `apidoc:"root,elem,usage-key"`
//      // 其它字段 ...
//  }
// 其 apidoc 标签值与其它的标签值格式相同，但是只有第一和第三个值是真实有效果的，
// 另两个值会被忽略。
package token

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
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

// 这些常量对应 Base 中相关字段的名称
const (
	rangeName         = "Range"
	usageKeyName      = "UsageKey"
	elementTagName    = "XMLName"
	elementTagEndName = "XMLNameEnd"
)

// Base 每一个 XML 节点必须包含的内容
type Base struct {
	core.Range

	// 表示对当前元素的一个说明内容的翻译 ID
	UsageKey message.Reference `apidoc:"-"`

	// 表示标签名或是属性名
	XMLName String `apidoc:"-"`

	// 表示标签的结束名称
	//
	// 如果是属性或是自闭合的标签，此值为空。
	XMLNameEnd String `apidoc:"-"`
}

// Usage 返回该节点的说明内容
func (b *Base) Usage() string {
	return locale.Sprintf(b.UsageKey)
}
