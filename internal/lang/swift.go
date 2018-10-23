// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lang

// swift 嵌套风格的块注释。会忽略掉内嵌的注释块。
type swiftNestMCommentBlock struct {
	begin      string
	end        string
	beginRunes []byte
	endRunes   []byte
	level      int8
}

func newSwiftNestMCommentBlock(begin, end string) Blocker {
	return &swiftNestMCommentBlock{
		begin:      begin,
		end:        end,
		beginRunes: []byte(begin),
		endRunes:   []byte(end),
	}
}

func (b *swiftNestMCommentBlock) BeginFunc(l *Lexer) bool {
	if l.Match(b.begin) {
		b.level++
		return true
	}

	return false
}

func (b *swiftNestMCommentBlock) EndFunc(l *Lexer) ([][]byte, bool) {
	lines := make([][]byte, 0, 20)
	line := make([]byte, 0, 100)

LOOP:
	for {
		switch {
		case l.AtEOF():
			return nil, false
		case l.Match(b.end):
			b.level--
			if b.level == 0 {
				lines = append(lines, filterSymbols(line, b.begin))
				break LOOP
			}

			line = append(line, b.endRunes...)
			continue LOOP
		case l.Match(b.begin):
			b.level++
			line = append(line, b.beginRunes...)
			continue LOOP
		default:
			r := l.data[l.pos]
			l.pos++
			line = append(line, r)
			if r == '\n' {
				lines = append(lines, filterSymbols(line, b.begin))
				line = make([]byte, 0, 100)
			}
		}
	}

	return lines, true
}
