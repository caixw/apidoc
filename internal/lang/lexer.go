// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lang

import (
	"bytes"
	"unicode"
)

// 是对一个文本内容的包装，方便 blocker 等接口操作。
type lexer struct {
	blocks []Blocker
	data   []byte
	pos    int

	ln    int // 上次记录的行号
	lnPos int // 上次记录行号时所在的位置
}

// 是否已经在文件末尾。
func (l *lexer) atEOF() bool {
	return l.pos >= len(l.data)
}

// 接下来的 n 个字符是否匹配指定的字符串，
// 若匹配，则将指定移向该字符串这后，否则不作任何操作。
func (l *lexer) match(word string) bool {
	if l.atEOF() || (l.pos+len(word) > len(l.data)) { // 剩余字符没有 word 长，直接返回 false
		return false
	}

	if bytes.HasPrefix(l.data[l.pos:], []byte(word)) {
		l.pos += len(word)
		return true
	}

	return false
}

var newLine = []byte("\n")

func (l *lexer) lineNumber() int {
	if l.lnPos < l.pos {
		l.ln += bytes.Count(l.data[l.lnPos:l.pos], newLine)
		l.lnPos = l.pos
	}

	return l.ln
}

// 从当前位置往后查找，直到找到第一个与 blocks 中某个相匹配的，并返回该 blocker 。
func (l *lexer) block() Blocker {
	for {
		if l.atEOF() {
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
func (l *lexer) skipSpace() {
	for {
		if l.atEOF() {
			return
		}

		r := l.data[l.pos]
		if !unicode.IsSpace(rune(r)) || r == '\n' {
			return
		}
		l.pos++
	}
}
