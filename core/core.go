// SPDX-License-Identifier: MIT

// Package core 提供基础的核心功能
package core

import "fmt"

const (
	// Name 程序的正式名称
	Name = "apidoc"

	// RepoURL 源码仓库地址
	RepoURL = "https://github.com/caixw/apidoc"

	// OfficialURL 官网
	OfficialURL = "https://apidoc.tools"

	// XMLNamespace 定义 xml 命名空间的 URI
	XMLNamespace = "https://apidoc.tools/v6/XMLSchema"
)

// Ranger Range 实现的方法集
//
// 所有内嵌 Range 的对象都可以使用此接口判断是否内嵌 Range。
type Ranger interface {
	Contains(Position) bool
	R() Range
	Equal(Range) bool
}

// Block 最基本的代码单位
//
// 一般从注释提取的一个完整注释作为 Block 实例。
type Block struct {
	Location Location
	Data     []byte
}

// Position 用于描述字符在文件中的定位
//
// 兼容 LSP https://microsoft.github.io/language-server-protocol/specifications/specification-current/#position
type Position struct {
	Line      int `json:"line" apidoc:"-"`
	Character int `json:"character" apidoc:"-"`
}

// Range 用于描述文档中的一段范围
//
// 兼容 LSP https://microsoft.github.io/language-server-protocol/specifications/specification-current/#range
type Range struct {
	Start Position `json:"start" apidoc:"-"`
	End   Position `json:"end" apidoc:"-"`
}

// Location 用于描述一段内容的定位
//
// 兼容 LSP https://microsoft.github.io/language-server-protocol/specifications/specification-current/#location
type Location struct {
	URI   URI   `json:"uri" apidoc:"-"`
	Range Range `json:"range" apidoc:"-"`
}

// Equal 判断与 v 是否相同
func (p Position) Equal(v Position) bool {
	return p.Line == v.Line && p.Character == v.Character
}

// Equal 判断与 v 是否相同
func (r Range) Equal(v Range) bool {
	return r.Start.Equal(v.Start) && r.End.Equal(v.End)
}

// IsEmpty 表示 Range 表示的范围长度为空
func (r Range) IsEmpty() bool {
	return r.End == r.Start
}

// Contains 是否包含了 p 这个点
func (r Range) Contains(p Position) bool {
	s := r.Start
	e := r.End
	return (s.Line < p.Line || (s.Line == p.Line && s.Character <= p.Character)) &&
		(e.Line > p.Line || (e.Line == p.Line && e.Character >= p.Character))
}

// R 返回当前的范围
func (r Range) R() Range {
	return r
}

func (l Location) String() string {
	s := l.Range.Start
	e := l.Range.End
	return fmt.Sprintf("%s[%d:%d,%d:%d]", l.URI, s.Line, s.Character, e.Line, e.Character)
}

// Equal 判断与 v 是否相等
//
// 所有字段都相同即返回 true。
func (l Location) Equal(v Location) bool {
	return l.Range.Equal(v.Range) && l.URI == v.URI
}
