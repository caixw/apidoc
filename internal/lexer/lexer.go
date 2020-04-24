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
	core.Block
	lastIndex int

	// 分别表示当前和之前的定位，可以在某些可撤消的操作之前保存定位信息到 prev
	current, prev Position
}

// New 声明 Lexer 实例
func New(b core.Block) (*Lexer, error) {
	// 以下代码主要保证内容都是合法的 utf8 编码，
	// 这样后续的操作不用再判断每个 utf8.DecodeRune 的调用返回是否都正常。
	p := Position{Position: b.Location.Range.Start}
	for {
		r, size := utf8.DecodeRune(b.Data[p.Offset:])
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
		Block:     b,
		lastIndex: p.Offset - 1,
		current:   Position{Position: b.Location.Range.Start},
		prev:      Position{Position: b.Location.Range.Start},
	}, nil
}

// AtEOF 是否已经结束
func (l *Lexer) AtEOF() bool {
	return l.Position().Offset > l.lastIndex
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
		r, size := utf8.DecodeRune(l.Data[p.Offset:])
		if size == 0 || r != w {
			return false
		}
		p = p.add(r, size)
	}

	l.prev = l.current
	l.current = p
	return true
}

// Position 返回当前在 data 中的偏移量
func (l *Lexer) Position() Position {
	return l.current
}

// Move 移动当前的分析器的位置
//
// 执行此操作之后，Rollback 将失效，且 AtEOF 为 false
func (l *Lexer) Move(p Position) {
	l.current = p
	l.prev.Offset = -1
}

// Spaces 获取之后的所有空格，不包含换行符
//
// NOTE: 可回滚该操作
func (l *Lexer) Spaces(exclude rune) []byte {
	l.prev = l.current

	for {
		r, size := utf8.DecodeRune(l.Data[l.current.Offset:])
		if size == 0 {
			break
		}

		if r == exclude || !unicode.IsSpace(r) { // 碰到换行符或是非空字符，则中止
			break
		}

		l.current = l.current.add(r, size)
	}

	return l.Bytes(l.prev.Offset, l.current.Offset)
}

// DelimString 查找 delim 并返回到此字符的所有内容
//
// contain 表示是否包含 delim 本身，如果为 false，则返回内容不包含，且该字符串会退回至输入流中，等待下次被读取。
//
// NOTE: 可回滚此操作
func (l *Lexer) DelimString(delim string, contain bool) ([]byte, bool) {
	if len(delim) == 0 {
		panic("参数 delim 不能为空值")
	}

	start := l.current
	for {
		if l.AtEOF() {
			l.current = start
			return nil, false
		}

		if l.Match(delim) {
			if !contain {
				l.Rollback()
			}
			l.prev = start
			return l.Bytes(start.Offset, l.current.Offset), true
		}
		l.Next(1)
	}
}

// Delim 查找 delim 并返回到此字符的所有内容
//
// NOTE: 可回滚此操作
func (l *Lexer) Delim(delim rune, contain bool) ([]byte, bool) {
	return l.DelimFunc(func(r rune) bool { return r == delim }, contain)
}

// DelimFunc 查找并返回当前位置到 f 确定位置的所有内容
//
// contain 表示是否包含字符本身，如果为 false，则返回内容不包含，且该字符会退回至输入流中，等待下次被读取。
//
// NOTE: 可回滚此操作
func (l *Lexer) DelimFunc(f func(r rune) bool, contain bool) ([]byte, bool) {
	if f == nil {
		panic("参数 f 不能为空")
	}

	p := l.current
	found := false

	for {
		r, size := utf8.DecodeRune(l.Data[p.Offset:])
		if size == 0 {
			break
		}
		p = p.add(r, size)

		if f(r) {
			found = true

			if !contain {
				p = p.sub(r, size)
			}
			break
		}
	}

	if !found {
		return nil, false
	}

	l.prev = l.current
	l.current = p
	return l.Bytes(l.prev.Offset, l.current.Offset), true
}

// Next 返回之后的 n 个字符，或是直到内容结束
//
// NOTE: 可回滚该操作
func (l *Lexer) Next(n int) []byte {
	l.prev = l.current

	for i := 0; i < n; i++ {
		r, size := utf8.DecodeRune(l.Data[l.Position().Offset:])
		if size == 0 {
			break
		}
		l.current = l.current.add(r, size)
	}

	return l.Bytes(l.prev.Offset, l.current.Offset)
}

// All 获取当前定位之后的所有内容
//
// NOTE: 可回滚该操作
func (l *Lexer) All() []byte {
	l.prev = l.current

	for {
		r, size := utf8.DecodeRune(l.Data[l.current.Offset:])
		if size == 0 {
			break
		}
		l.current = l.current.add(r, size)
	}

	return l.Data[l.prev.Offset:]
}

// Rollback 回滚上一次的操作
func (l *Lexer) Rollback() {
	if l.prev.Offset != -1 {
		l.current = l.prev // 回滚
		l.prev.Offset = -1 // 清空 prev
	}
}

// Bytes 返回指定范围的内容
//
// NOTE: 并不会改变定位信息
func (l *Lexer) Bytes(start, end int) []byte {
	return l.Data[start:end]
}
