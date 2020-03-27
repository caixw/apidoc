// SPDX-License-Identifier: MIT

// Package lexer 提供基本的分词功能
package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
)

// Lexer 是对一个文本内容的包装，方便 blocker 等接口操作。
type Lexer struct {
	data  []byte
	atEOF bool

	// 分别表示当前和之前的定位，可以在某些可撤消的操作之前保存定位信息到 prev
	current, prev Position
}

// New 声明 Lexer 实例
func New(data []byte) (*Lexer, error) {
	// 以下代码主要保证内容都是合法的 utf8 编码，
	// 这样后续的操作不用再判断每个 utf8.DecodeRune 的调用返回是否都正常。
	p := Position{}
	for {
		r, size := utf8.DecodeRune(data[p.Offset:])
		if size == 0 {
			break
		}
		if r == utf8.RuneError && size == 1 {
			loc := core.Location{
				Range: core.Range{
					Start: p.Position,
					End:   core.Position{Line: p.Line, Character: p.Character + size},
				},
			}
			return nil, core.NewLocaleError(loc, "", locale.ErrInvalidUTF8Character)
		}

		p.Offset += size
		p.Character++
		if r == '\n' {
			p.Line++
			p.Character = 0
		}
	}

	return &Lexer{
		data: data,
	}, nil
}

// AtEOF 是否已经结束
func (l *Lexer) AtEOF() bool {
	return l.atEOF
}

// Match 接下来的 n 个字符是否匹配指定的字符串，
// 若匹配，则将指定移向该字符串这后，否则不作任何操作。
//
// NOTE: 可回滚该操作
func (l *Lexer) Match(word string) bool {
	if word == "" {
		return false
	}

	p := l.current

	for _, w := range word {
		r, size := utf8.DecodeRune(l.data[p.Offset:])
		if size == 0 || r != w {
			return false
		}

		p.Offset += size
		p.Character++
		if r == '\n' {
			p.Line++
			p.Character = 0
		}
	}

	l.prev = l.current
	l.current = p
	return true
}

// Position 返回当前在 data 中的偏移量
func (l *Lexer) Position() Position {
	return l.current
}

// Spaces 获取之后的所有空格，不包含换行符
//
// NOTE: 可回滚该操作
func (l *Lexer) Spaces() []byte {
	p := l.current

	for {
		r, size := utf8.DecodeRune(l.data[p.Offset:])
		if size == 0 {
			l.atEOF = true
			break
		}

		if r == '\n' || !unicode.IsSpace(r) { // 碰到换行符或是非空字符，则中止
			break
		}

		p.Offset += size
		p.Character++
	}

	l.prev = l.current
	l.current = p
	return l.Bytes(l.prev.Offset, l.current.Offset)
}

// DelimString 查找 delim 并返回到此字符的所有内容，未找到则返回空值
//
// NOTE: 可回滚此操作
func (l *Lexer) DelimString(delim string) []byte {
	if len(delim) == 0 {
		return nil
	}

	delimRunes := []rune(delim)
	if len(delimRunes) == 1 {
		return l.Delim(delimRunes[0])
	}

	curr := l.current
	prev := l.prev
	for {
		if l.AtEOF() { // 一直到结束都未找到匹配项，则还原到起始位置
			l.current = curr
			l.prev = prev
			return nil
		}

		// 找到第一个匹配的字符
		if l.Delim(delimRunes[0]) == nil {
			return nil
		}

		// 查看后续字符是否匹配，如果不匹配，则继续从当前位置开始查找
		if l.Match(string(delimRunes[1:])) {
			return l.Bytes(curr.Offset, l.current.Offset)
		}
	}
}

// Delim 查找 delim 并返回到此字符的所有内容，未找到则返回空值
//
// NOTE: 可回滚此操作
func (l *Lexer) Delim(delim rune) []byte {
	if delim == 0 {
		return nil
	}

	p := l.current
	found := false

	for {
		r, size := utf8.DecodeRune(l.data[p.Offset:])
		if size == 0 {
			l.atEOF = true
			break
		}

		p.Offset += size
		p.Character++
		if r == '\n' {
			p.Line++
			p.Character = 0
		}

		if r == delim {
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	l.prev = l.current
	l.current = p
	return l.Bytes(l.prev.Offset, l.current.Offset)
}

// Next 返回之后的 n 个字符，或是直到内容结束
//
// NOTE: 可回滚该操作
func (l *Lexer) Next(n int) []byte {
	p := l.current

	for i := 0; i < n; i++ {
		r, size := utf8.DecodeRune(l.data[p.Offset:])
		if size == 0 {
			l.atEOF = true
			break
		}

		p.Offset += size
		p.Character++
		if r == '\n' {
			p.Line++
			p.Character = 0
		}
	}

	l.prev = l.current
	l.current = p
	return l.Bytes(l.prev.Offset, l.current.Offset)
}

// All 获取当前定位之后的所有内容
//
// NOTE: 可回滚该操作
func (l *Lexer) All() []byte {
	p := l.current

	for {
		r, size := utf8.DecodeRune(l.data[p.Offset:])
		if size == 0 {
			l.atEOF = true
			break
		}

		p.Offset += size
		p.Character++
		if r == '\n' {
			p.Line++
			p.Character = 0
		}
	}

	l.prev = l.current
	l.current = p
	return l.data[l.prev.Offset:]
}

// Rollback 回滚操作
func (l *Lexer) Rollback() {
	if l.prev.Offset == 0 {
		return
	}

	// 回滚
	l.current = l.prev
	l.atEOF = false

	l.prev.Offset = 0 // 清空 prev
}

// Bytes 返回指定范围的内容
func (l *Lexer) Bytes(start, end int) []byte {
	return l.data[start:end]
}
