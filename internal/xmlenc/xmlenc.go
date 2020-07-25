// SPDX-License-Identifier: MIT

// Package xmlenc 有关文档格式编解码处理
package xmlenc

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Base 所有文档节点的基本元素
type Base struct {
	core.Range
	UsageKey message.Reference `apidoc:"-"` // 表示对当前元素的一个说明内容的翻译 ID
}

// BaseAttribute 所有 XML 属性节点的基本元素
type BaseAttribute struct {
	Base
	AttributeName Name `apidoc:"-"`
}

// BaseTag 所有 XML 标签的基本元素
type BaseTag struct {
	Base
	StartTag Name `apidoc:"-"` // 表示起始标签名
	EndTag   Name `apidoc:"-"` // 表示标签的结束名称，如果是自闭合的标签，此值为空。
}

// Usage 本地化的当前字段介绍内容
func (b *Base) Usage() string {
	return locale.Sprintf(b.UsageKey)
}
