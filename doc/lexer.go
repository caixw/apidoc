// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"unicode"

	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/input"
)

// 简单的词法分析
type lexer struct {
	data      input.Block
	pos       int       // 当前指针位置
	ln        int       // 当前位置所在的行号
	backupTag *lexerTag // 回退的标签内容
	h         *errors.Handler
}

// 表示一个标签的内容
type lexerTag struct {
	File string
	Line int // 当前 tag 在文件中的起始行号
	Data []byte
	Name string // 标签名称
	l    *lexer
}

func newLexer(block input.Block, h *errors.Handler) *lexer {
	return &lexer{
		data: block,
		ln:   block.Line,
		h:    h,
	}
}

func newLexerTag(l *lexer, file string, line int, data []byte) *lexerTag {
	strs := splitWords(data, 2)

	tag := &lexerTag{
		File: file,
		Line: line,
		Name: string(strs[0]),
		l:    l,
	}

	if len(strs) == 2 {
		tag.Data = strs[1]

		// 去掉最后的换行符
		lastIndex := len(tag.Data) - 1
		if tag.Data[lastIndex] == '\n' {
			tag.Data = tag.Data[:lastIndex]
		}
	}

	return tag
}

func (l *lexer) atEOF() bool {
	return l.pos >= len(l.data.Data)
}

// Backup 回退一个标签
func (l *lexer) backup(t *lexerTag) {
	l.backupTag = t
}

// Tag 返回下一个标签。
//
// 若返回 nil 表示已经在结尾处。
func (l *lexer) tag() (t *lexerTag) {
	if l.backupTag != nil {
		t = l.backupTag
		l.backupTag = nil
		return t
	}

	newLine := false
	start := l.pos
	end := l.pos
	ln := l.ln

LOOP:
	for ; ; l.pos++ {
		if l.atEOF() {
			data := l.data.Data[start:l.pos]
			if len(data) == 0 {
				return nil
			}
			return newLexerTag(l, l.data.File, ln, data)
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
			return newLexerTag(l, l.data.File, ln, l.data.Data[start:end])
		default:
			newLine = false
		}
	}
}

// 将 tag.Data 以空隔分成 num 个数组。
//
// 如果不够数量，则返回实际数量的元素。
func (tag *lexerTag) words(num int) [][]byte {
	return splitWords(tag.Data, num)
}

// 分成指定的行并返回
//
// 如果不够数量，则返回实际数量的元素。
func (tag *lexerTag) lines(num int) [][]byte {
	return splitLines(tag.Data, num)
}

// 将 data 以空隔分成 num 个数组。
//
// 如果不够数量，则返回实际数量的元素。
func splitWords(data []byte, size int) [][]byte {
	return splitFunc(data, size, func(b byte) bool { return unicode.IsSpace(rune(b)) })
}

func splitLines(data []byte, size int) [][]byte {
	return splitFunc(data, size, func(b byte) bool { return b == '\n' })
}

func splitFunc(data []byte, size int, fn func(b byte) bool) [][]byte {
	ret := make([][]byte, 0, size)
	start := 0
	pos := 0
	issperator := true // 前一个字符是否为分隔符

	for ; ; pos++ {
		switch {
		case pos >= len(data): // EOF
			if !issperator { // 如果依然为 true，说明剩余的都是分隔符，不返回内容
				ret = append(ret, data[start:])
			}
			return ret
		case fn(data[pos]):
			if !issperator {
				ret = append(ret, data[start:pos])
				start = pos
				issperator = true
			}
		default:
			if issperator {
				if len(ret) >= size-1 {
					return append(ret, data[pos:])
				}

				start = pos
				issperator = false
			}
		} // end switch
	} // end for
}
