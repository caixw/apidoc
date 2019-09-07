// SPDX-License-Identifier: MIT

package lang

// swift 嵌套风格的块注释。会忽略掉内嵌的注释块。
type swiftNestMCommentBlock struct {
	begin      string
	end        string
	prefix     string // 需要过滤的前缀
	beginRunes []byte
	endRunes   []byte
	level      int8
}

func newSwiftNestMCommentBlock(begin, end, prefix string) Blocker {
	return &swiftNestMCommentBlock{
		begin:      begin,
		end:        end,
		prefix:     prefix,
		beginRunes: []byte(begin),
		endRunes:   []byte(end),
	}
}

func (b *swiftNestMCommentBlock) BeginFunc(l *lexer) bool {
	if l.match(b.begin) {
		b.level++
		return true
	}

	return false
}

func (b *swiftNestMCommentBlock) EndFunc(l *lexer) ([][]byte, bool) {
	lines := make([][]byte, 0, 20)
	line := make([]byte, 0, 100)

LOOP:
	for {
		switch {
		case l.atEOF():
			return nil, false
		case l.match(b.end):
			b.level--
			if b.level == 0 {
				if len(line) > 0 { // 如果 len(line) == 0 表示最后一行仅仅只有一个结束符
					lines = append(lines, filterSymbols(line, b.prefix))
				}
				break LOOP
			}

			line = append(line, b.endRunes...)
			continue LOOP
		case l.match(b.begin):
			b.level++
			line = append(line, b.beginRunes...)
			continue LOOP
		default:
			r := l.data[l.pos]
			l.pos++
			line = append(line, r)
			if r == '\n' {
				lines = append(lines, filterSymbols(line, b.prefix))
				line = make([]byte, 0, 100)
			}
		}
	}

	return lines, true
}
