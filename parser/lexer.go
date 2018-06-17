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
	file string
	ln   int // 当前 tag 在文件中的起始行号
	data []byte
	name []byte // 标签名称
}

// 声明一个新的 lexer 实例。
func newLexer(block input.Block) *lexer {
	return &lexer{
		data: block,
		ln:   block.Line,
	}
}

func newTag(file string, line int, data []byte) *tag {
	strs := split(data, 2)

	tag := &tag{
		file: file,
		ln:   line,
		name: strs[0],
	}

	if len(strs) == 2 {
		tag.data = strs[1]
	}

	return tag
}

func (l *lexer) atEOF() bool {
	return l.pos >= len(l.data.Data)
}

func (l *lexer) tag() (t *tag, eof bool) {
	newLine := false
	start := l.pos
	end := l.pos
	ln := l.ln

LOOP:
	for ; ; l.pos++ {
		if l.atEOF() {
			data := l.data.Data[start:l.pos]
			if len(data) == 0 {
				return nil, true
			}
			return newTag(l.data.File, ln, data), true
		}

		b := l.data.Data[l.pos]
		switch {
		case b == '\n':
			l.ln++
			newLine = true
			end = l.pos + 1 // 包含当前的换行符
		case newLine && unicode.IsSpace(rune(b)): // 跳过行首空白字符
			continue LOOP
		case newLine && b == '@':
			return newTag(l.data.File, ln, l.data.Data[start:end]), false
		default:
			newLine = false
		}
	}
}

func (t *tag) syntaxError(message string) error {
	return syntaxError(message, t.file, t.ln)
}

func split(data []byte, size int) [][]byte {
	ret := make([][]byte, 0, size)
	start := 0
	pos := 0
	isspace := true // 前一个字符是否为空白字符

	for ; ; pos++ {
		switch {
		case pos >= len(data): // EOF
			return append(ret, data[start:])
		case unicode.IsSpace(rune(data[pos])):
			if !isspace {
				ret = append(ret, data[start:pos])
				start = pos
				isspace = true
			}
		default:
			if isspace {
				if len(ret) >= size-1 {
					return append(ret, data[pos:])
				}

				start = pos
				isspace = false
			}
		}
	}
}
