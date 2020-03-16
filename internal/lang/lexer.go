// SPDX-License-Identifier: MIT

package lang

import (
	"unicode"
	"unicode/utf8"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
)

// Error 表示解析错误
type Error struct {
	Position core.Position
	Message  string
}

func (err *Error) Error() string {
	return err.Message
}

type position struct {
	core.Position

	// 表示的是字节的偏移量，
	// 而 Position.Character 表示的是当前行`字符`的偏移量
	Offset int
}

// Lexer 是对一个文本内容的包装，方便 blocker 等接口操作。
type Lexer struct {
	blocks []Blocker
	data   []byte
	atEOF  bool

	// 分别表示当前和之前的定位，可以在某些可撤消的操作之前保存定位信息到 prev
	current, prev position
}

// NewLexer 声明 Lexer 实例
func NewLexer(data []byte, blocks []Blocker) (*Lexer, error) {
	// 以下代码主要保证内容都是合法的 utf8 编码，
	// 这样后续的操作不用再判断每个 utf8.DecodeRune 的调用返回是否都正常。
	p := position{}
	for {
		r, size := utf8.DecodeRune(data[p.Offset:])
		if size == 0 {
			break
		}
		if r == utf8.RuneError {
			return nil, &Error{Position: p.Position, Message: locale.Sprintf(locale.ErrInvalidUTF8Character)}
		}

		p.Offset += size
		p.Character++
		if r == '\n' {
			p.Line++
			p.Character = 0
		}
	}

	return &Lexer{
		data:   data,
		blocks: blocks,
	}, nil
}

// AtEOF 是否已经在文件末尾。
func (l *Lexer) AtEOF() bool {
	return l.atEOF
}

// 接下来的 n 个字符是否匹配指定的字符串，
// 若匹配，则将指定移向该字符串这后，否则不作任何操作。
//
// NOTE: 可回滚该操作
func (l *Lexer) match(word string) bool {
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

// Position 当前在 data 中的偏移量
func (l *Lexer) Position() core.Position {
	return l.current.Position
}

// Block 从当前位置往后查找，直到找到第一个与 blocks 中某个相匹配的，并返回该 Blocker 。
func (l *Lexer) Block() (Blocker, core.Position) {
	for {
		if l.AtEOF() {
			return nil, core.Position{}
		}

		pos := l.Position()
		for _, block := range l.blocks {
			if block.BeginFunc(l) {
				return block, pos
			}
		}

		l.next(1)
	}
}

// 跳过之后的空格，不包含换行符
//
// NOTE: 可回滚该操作
func (l *Lexer) skipSpace() []byte {
	p := l.current
	bs := make([]byte, 0, 50)

	for {
		r, size := utf8.DecodeRune(l.data[p.Offset:])
		if size == 0 {
			l.atEOF = true
			break
		}

		if r == '\n' || !unicode.IsSpace(r) { // 碰到换行符或是非空字符，则中止
			break
		}

		bs = append(bs, l.data[p.Offset:p.Offset+size]...)
		p.Offset += size
		p.Character++
	}

	l.prev = l.current
	l.current = p
	return bs
}

// delim 查找 delim 并返回到此字符的所有内容，未找到则返回空值
//
// NOTE: 可回滚此操作
func (l *Lexer) delim(delim rune) []byte {
	p := l.current
	bs := make([]byte, 0, 50)
	found := false

	for {
		r, size := utf8.DecodeRune(l.data[p.Offset:])
		if size == 0 {
			l.atEOF = true
			break
		}
		bs = append(bs, l.data[p.Offset:p.Offset+size]...)

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
	return bs
}

// 返回之后的 n 个字符，或是直到内容结束
//
// NOTE: 可回滚该操作
func (l *Lexer) next(n int) []byte {
	p := l.current
	bs := make([]byte, 0, n)

	for i := 0; i < n; i++ {
		r, size := utf8.DecodeRune(l.data[p.Offset:])
		if size == 0 {
			l.atEOF = true
			break
		}
		bs = append(bs, l.data[p.Offset:p.Offset+size]...)

		p.Offset += size
		p.Character++
		if r == '\n' {
			p.Line++
			p.Character = 0
		}
	}

	l.prev = l.current
	l.current = p
	return bs
}

// 获取当前定位之后的所有内容
//
// NOTE: 可回滚该操作
func (l *Lexer) all() []byte {
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

// 回滚操作
func (l *Lexer) back() {
	if l.prev.Offset == 0 {
		return
	}

	// 回滚
	l.current = l.prev
	l.atEOF = false

	l.prev.Offset = 0 // 清空 prev
}
