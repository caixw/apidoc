// SPDX-License-Identifier: MIT

package lang

import (
	"bytes"
	"unicode"
)

// Lexer 是对一个文本内容的包装，方便 blocker 等接口操作。
type Lexer struct {
	blocks []Blocker
	data   []byte
	pos    int

	ln    int // 上次记录的行号
	lnPos int // 上次记录行号时所在的位置
}

// NewLexer 声明 Lexer 实例
func NewLexer(data []byte, blocks []Blocker) *Lexer {
	return &Lexer{
		data:   data,
		blocks: blocks,
	}
}

// AtEOF 是否已经在文件末尾。
func (l *Lexer) AtEOF() bool {
	return l.pos >= len(l.data)
}

// 接下来的 n 个字符是否匹配指定的字符串，
// 若匹配，则将指定移向该字符串这后，否则不作任何操作。
func (l *Lexer) match(word string) bool {
	if l.AtEOF() || (l.pos+len(word) > len(l.data)) { // 剩余字符没有 word 长，直接返回 false
		return false
	}

	if bytes.HasPrefix(l.data[l.pos:], []byte(word)) {
		l.pos += len(word)
		return true
	}

	return false
}

var newLine = []byte("\n")

// LineNumber 获取当前位置所在的行号
func (l *Lexer) LineNumber() int {
	if l.lnPos < l.pos {
		l.ln += bytes.Count(l.data[l.lnPos:l.pos], newLine)
		l.lnPos = l.pos
	}

	return l.ln
}

// Block 从当前位置往后查找，直到找到第一个与 blocks 中某个相匹配的，并返回该 Blocker 。
func (l *Lexer) Block() Blocker {
	for {
		if l.AtEOF() {
			return nil
		}

		for _, block := range l.blocks {
			if block.BeginFunc(l) {
				return block
			}
		}

		l.pos++
	}
}

// 跳过除换行符以外的所有空白字符。
func (l *Lexer) skipSpace() {
	for {
		if l.AtEOF() {
			return
		}

		r := l.data[l.pos]
		if !unicode.IsSpace(rune(r)) || r == '\n' {
			return
		}
		l.pos++
	}
}

// 读取到当前行行尾。
//
// 返回 nil 表示没有换行符，即当前就是最后一行。
func (l *Lexer) line() []byte {
	start := l.pos
	for index, b := range l.data[l.pos:] {
		if l.AtEOF() {
			return nil
		}

		if b == '\n' {
			l.pos += index
			return l.data[start : index+start]
		}
	} // end for

	return nil
}
