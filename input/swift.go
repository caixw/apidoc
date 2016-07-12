// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

// swift 嵌套风格的块注释。会忽略掉内嵌的注释块。
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
				lines = append(lines, filterSymbols(line, b.begin))
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
				lines = append(lines, filterSymbols(line, b.begin))
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
