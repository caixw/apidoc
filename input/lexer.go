// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"bytes"
	"unicode/utf8"

	"github.com/caixw/apidoc/app"
)

// lexer 是对一个文本内容的包装，方便 blocker 等接口操作。
type lexer struct {
	data    []byte
	pos     int
	isAtEOF bool

	ln    int // 上次记录的行号
	lnPos int // 上次记录行号时所在的位置
}

// 是否已经在文件末尾。
func (l *lexer) atEOF() bool {
	return l.isAtEOF || l.pos >= len(l.data)
}

// 获取当前的字符，并将指针指向下一个字符。
func (l *lexer) next() rune {
	r, w := utf8.DecodeRune(l.data[l.pos:])
	l.pos += w

	if r == utf8.RuneError && w == 0 { // EOF
		l.isAtEOF = true
	}

	return r
}

// 接下来的 n 个字符是否匹配指定的字符串，
// 若匹配，则将指定移向该字符串这后，否则不作任何操作。
func (l *lexer) match(word string) bool {
	if l.atEOF() || (l.pos+len(word) > len(l.data)) { // 剩余字符没有 word 长，直接返回 false
		return false
	}

	width := 0
	for _, r := range word {
		rr, w := utf8.DecodeRune(l.data[l.pos:])
		if rr != r {
			l.pos -= width
			return false
		}

		l.pos += w
		width += w
	}

	return true
}

var newLine = []byte("\n")

func (l *lexer) lineNumber() int {
	if l.lnPos < l.pos {
		l.ln += bytes.Count(l.data[l.lnPos:l.pos], newLine)
		l.lnPos = l.pos
	}

	return l.ln
}

// 构建一条语法错误的信息。
func (l *lexer) syntaxError(msg string) *app.SyntaxError {
	return &app.SyntaxError{
		Line:    l.lineNumber(),
		Message: msg,
	}
}

// 从当前位置往后查找，直到找到第一个与 blocks 中某个相匹配的，并返回该 blocker 。
func (l *lexer) block(blocks []blocker) blocker {
	for {
		if l.atEOF() {
			return nil
		}

		for _, block := range blocks {
			if block.BeginFunc(l) {
				return block
			}
		}

		l.next()
	}
}
