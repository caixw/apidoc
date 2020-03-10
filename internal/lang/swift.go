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
	raw = make([]byte, 0, 200)
	line := make([]byte, 0, 100)

LOOP:
	for {
		switch {
		case l.AtEOF():
			return nil, nil, false
		case l.match(b.end):
			b.level--
			if b.level == 0 {
				if len(line) > 0 { // 如果 len(line) == 0 表示最后一行仅仅只有一个结束符
					data = append(data, filterSymbols(line, b.prefix)...)
				}
				break LOOP
			}

			raw = append(raw, b.ends...)
			line = append(line, b.ends...)
			continue LOOP
		case l.match(b.begin):
			if b.level > 0 {
				raw = append(raw, b.begins...)
			}

			b.level++
			line = append(line, b.begins...)
			continue LOOP
		default:
			r := l.data[l.offset]
			raw = append(raw, r)
			l.offset++
			line = append(line, r)
			if r == '\n' {
				data = append(data, filterSymbols(line, b.prefix)...)
				line = make([]byte, 0, 100)
			}
		}
	}

	return raw, data, true
}
