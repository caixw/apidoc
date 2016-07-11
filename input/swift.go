// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"strings"
	"unicode"
)

type swiftNestMCommentBlock struct {
	begin      string
	end        string
	beginRunes []rune
	endRunes   []rune
	level      int8
}

func newSwiftNestMCommentBlock(begin, end string) blocker {
	return &swiftNestMCommentBlock{
		begin:      begin,
		end:        end,
		beginRunes: []rune(begin),
		endRunes:   []rune(end),
	}
}

func (b *swiftNestMCommentBlock) BeginFunc(l *lexer) bool {
	if l.match(b.begin) {
		b.level++
		return true
	}

	return false
}

func (b *swiftNestMCommentBlock) EndFunc(l *lexer) ([]rune, bool) {
	lines := make([][]rune, 0, 20)
	line := make([]rune, 0, 100)

LOOP:
	for {
		switch {
		case l.atEOF():
			return nil, false
		case l.match(b.end):
			b.level--
			if b.level == 0 {
				lines = append(lines, b.filterSymbols(line))
				break LOOP
			}

			line = append(line, b.endRunes...)
			continue LOOP
		case l.match(b.begin):
			b.level++
			line = append(line, b.beginRunes...)
			continue LOOP
		default:
			r := l.next()
			line = append(line, r)
			if r == '\n' {
				lines = append(lines, b.filterSymbols(line))
				line = make([]rune, 0, 100)
			}
		}
	}

	ret := make([]rune, 0, 1000)
	for _, v := range lines {
		ret = append(ret, v...)
	}
	return ret, true
}

// 行首若出现`空白字符+symbol+空白字符`的组合，则去掉这些字符。
// symbol 为 b.Begin 中的任意一个字符。
func (b *swiftNestMCommentBlock) filterSymbols(line []rune) []rune {
	for k, v := range line {
		if unicode.IsSpace(v) { // 跳过行首的空格
			continue
		}

		// 不存在指定的符号，直接返回原数据
		if strings.IndexRune(b.begin, v) < 0 {
			return line
		}

		// 若下个字符正好是是空格
		if len(line) > k+1 && unicode.IsSpace(line[k+1]) {
			return line[k+2:]
		} else {
			return line
		}
	}

	return line
}
