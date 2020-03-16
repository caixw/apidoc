// SPDX-License-Identifier: MIT

package lang

// swift 嵌套风格的块注释。会忽略掉内嵌的注释块。
type swiftNestMCommentBlock struct {
	begin  string
	end    string
	prefix []byte // 需要过滤的前缀
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
		prefix: []byte(prefix),
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
	raw = append(make([]byte, 0, 200), b.begins...)

LOOP:
	for {
		switch {
		case l.AtEOF():
			return nil, nil, false
		case l.match(b.end):
			raw = append(raw, b.ends...)
			b.level--
			if b.level == 0 {
				break LOOP
			}
		case l.match(b.begin):
			raw = append(raw, b.begins...)
			b.level++
		default:
			raw = append(raw, l.next(1)...)
		}
	}

	return raw, convertMultipleCommentToXML(raw, b.begins, b.ends, b.prefix), true
}
