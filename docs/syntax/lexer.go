// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package syntax

import (
	"unicode"

	"github.com/caixw/apidoc/input"
)

// Lexer 简单的词法分析
type Lexer struct {
	data      input.Block
	pos       int  // 当前指针位置
	ln        int  // 当前位置所在的行号
	backupTag *Tag // 回退的标签内容
}

// Tag 表示一个标签的内容
type Tag struct {
	File string
	Line int // 当前 tag 在文件中的起始行号
	Data []byte
	Name []byte // 标签名称
}

// NewLexer 声明一个新的 Lexer 实例。
func NewLexer(block input.Block) *Lexer {
	return &Lexer{
		data: block,
		ln:   block.Line,
	}
}

func newTag(file string, line int, data []byte) *Tag {
	strs := split(data, 2)

	tag := &Tag{
		File: file,
		Line: line,
		Name: strs[0],
	}

	if len(strs) == 2 {
		tag.Data = strs[1]
	}

	return tag
}

func (l *Lexer) atEOF() bool {
	return l.pos >= len(l.data.Data)
}

// Backup 回退一个标签
func (l *Lexer) Backup(t *Tag) {
	l.backupTag = t
}

// Tag 返回下一个标签
// eof 表示是否已经是结尾处。
func (l *Lexer) Tag() (t *Tag, eof bool) {
	if l.backupTag != nil {
		t = l.backupTag
		l.backupTag = nil
		return t, l.atEOF()
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

// Split 将 tag.Data 以空隔分成 num 个数组。
// 如果不够数量，则返回实际数量的元素。
func (tag *Tag) Split(num int) [][]byte {
	return split(tag.Data, num)
}
