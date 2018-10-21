// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

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
	Name string // 标签名称
}

// New 声明一个新的 Lexer 实例。
func New(block input.Block) *Lexer {
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
		Name: string(strs[0]),
	}

	if len(strs) == 2 {
		tag.Data = strs[1]
	}

	// 去掉最后的换行符
	lastIndex := len(tag.Data) - 1
	if tag.Data[lastIndex] == '\n' {
		tag.Data = tag.Data[:lastIndex]
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

// 分隔成指定大小的字符串数组
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
		} // end switch
	} // end for
}
