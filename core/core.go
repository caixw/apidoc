// SPDX-License-Identifier: MIT

// Package core 一些公共处理功能
package core

import "fmt"

// Block 表示原始的注释代码块
type Block struct {
	Location Location

	// Raw 表示原始的注释代码内容
	//
	// Data 为处理之后的数据
	// 为一个正常的 XML 格式内容，且长度应该与 Raw 相同。
	Raw  []byte
	Data []byte
}

// Position 用于描述字符在文件中的定位
//
// 兼容 LSP https://microsoft.github.io/language-server-protocol/specifications/specification-current/#position
type Position struct {
	Line      int `json:"line" apidoc:"-"`
	Character int `json:"character" apidoc:"-"`
}

// Range 用于描述文档在文件中的范围
//
// 兼容 LSP https://microsoft.github.io/language-server-protocol/specifications/specification-current/#range
type Range struct {
	Start Position `json:"start" apidoc:"-"`
	End   Position `json:"end" apidoc:"-"`
}

// Location 用于描述文档的定位
//
// 兼容 LSP https://microsoft.github.io/language-server-protocol/specifications/specification-current/#location
type Location struct {
	URI   URI   `json:"uri" apidoc:"-"`
	Range Range `json:"range" apidoc:"-"`
}

// Equal 判断两个 Position 是否相等
func (p Position) Equal(v Position) bool {
	return p.Line == v.Line && p.Character == v.Character
}

// IsEmpty 表示 Range 表示的范围长度为空
func (r Range) IsEmpty() bool {
	return r.End == r.Start
}

func (l Location) String() string {
	s := l.Range.Start
	e := l.Range.End
	return fmt.Sprintf("%s[%d:%d,%d:%d]", l.URI, s.Line, s.Character, e.Line, e.Character)
}
