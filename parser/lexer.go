// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"unicode"

	"github.com/caixw/apidoc/input"
)

// 简单的词法分析
type lexer struct {
	data input.Block
	pos  int // 当前指针位置
	ln   int // 当前位置所在的行号
}

// 表示一个标签的内容
type tag struct {
	ln   int // 当前 tag 在文件中的起始行号
	data []byte
}

// 声明一个新的 lexer 实例。
func newLexer(block input.Block) *lexer {
	return &lexer{
		data: block,
		ln:   block.Line,
	}
}

func (l *lexer) atEOF() bool {
	return l.pos >= len(l.data.Data)
}

func (l *lexer) tag() *tag {
	newLine := false
	start := l.pos

LOOP:
	for ; ; l.pos++ {
		if l.atEOF() {
			return &tag{data: l.data.Data[start:l.pos], ln: l.ln}
		}

		b := l.data.Data[l.pos]
		switch {
		case b == '\n':
			l.ln++
			newLine = true
		case newLine && unicode.IsSpace(rune(b)): // 跳过行首空白字符
			continue LOOP
		case newLine && b == '@':
			return &tag{data: l.data.Data[start:l.pos], ln: l.ln}
		default:
			newLine = false
		}
	}
}

func (t *tag) split(size int) [][]byte {
	ret := make([][]byte, 0, size)
	start := 0
	pos := 0
	isspace := true // 前一个字符是否为空白字符

	for ; ; pos++ {
		switch {
		case pos >= len(t.data): // EOF
			return append(ret, t.data[start:])
		case unicode.IsSpace(rune(t.data[pos])):
			if !isspace {
				ret = append(ret, t.data[start:pos])
				start = pos
				isspace = true
			}
		default:
			if isspace {
				if len(ret) >= size-1 {
					return append(ret, t.data[pos:])
				}

				start = pos
				isspace = false
			}
		}
	}
}
