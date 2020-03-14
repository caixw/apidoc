// SPDX-License-Identifier: MIT

package lang

// swift 嵌套风格的块注释。会忽略掉内嵌的注释块。
type swiftNestMCommentBlock struct {
	begin  string
	end    string
	prefix string // 需要过滤的前缀
	begins []byte
	ends   []byte
	level  int8
}

// prefix 表示每一行的前缀符号，比如：
//  /*
//   *
//   */
// 中的 * 字符
func newSwiftNestMCommentBlock(begin, end, prefix string) Blocker {
	return &swiftNestMCommentBlock{
		begin:  begin,
		end:    end,
		prefix: prefix,
		begins: []byte(begin),
		ends:   []byte(end),
	}
}

func (b *swiftNestMCommentBlock) BeginFunc(l *Lexer) bool {
	if l.match(b.begin) {
		b.level++
		return true
	}

	return false
}

func (b *swiftNestMCommentBlock) EndFunc(l *Lexer) (raw, data []byte, ok bool) {
	data = make([]byte, 0, 200)
	raw = append(make([]byte, 0, 200), b.begins...)
	line := make([]byte, 0, 100)

LOOP:
	for {
		switch {
		case l.AtEOF():
			return nil, nil, false
		case l.match(b.end):
			raw = append(raw, b.ends...)

			b.level--
			if b.level == 0 {
				if len(line) > 0 { // 如果 len(line) == 0 表示最后一行仅仅只有一个结束符
					data = append(data, filterSymbols(line, b.prefix)...)
					line = line[:0]
				}
				break LOOP
			}

			line = append(line, b.ends...)
		case l.match(b.begin):
			b.level++
			raw = append(raw, b.begins...)
			line = append(line, b.begins...)
		default:
			bs := l.next(1)
			raw = append(raw, bs...)
			line = append(line, bs...)
			if isNewline(bs) {
				data = append(data, filterSymbols(line, b.prefix)...)
				line = line[:0]
			}
		}
	}

	return raw, data, true
}
