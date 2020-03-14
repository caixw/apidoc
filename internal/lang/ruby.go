// SPDX-License-Identifier: MIT

package lang

// 表示超始和结束符号必须占满一行的情况
type rubyMultipleComment struct {
	begin, end, escape string
	begins, ends       []byte
}

func newRubyMultipleComment(begin, end, escape string) Blocker {
	begin += "\n"
	end += "\n"
	return &rubyMultipleComment{
		begin:  begin,
		end:    end,
		escape: escape,
		begins: []byte(begin),
		ends:   []byte(end),
	}
}

// BeginFunc 实现 Blocker.BeginFunc
func (b *rubyMultipleComment) BeginFunc(l *Lexer) bool {
	return l.Position().Character == 0 && l.match(b.begin)
}

// 从 l 的当前位置一直到定义的 b.End 之间的所有字符。
// 会对每一行应用 filterSymbols 规则。
func (b *rubyMultipleComment) EndFunc(l *Lexer) (raw, data []byte, ok bool) {
	data = make([]byte, 0, 200)
	raw = append(make([]byte, 0, 200), b.begins...)
	line := make([]byte, 0, 100)

	for {
		switch {
		case l.AtEOF(): // 没有找到结束符号，直接到达文件末尾
			return nil, nil, false
		case l.Position().Character == 0 && l.match(b.end):
			raw = append(raw, b.ends...)
			return raw, data, true
		default:
			bs := l.next(1)
			raw = append(raw, bs...)
			line = append(line, bs...)
			if isNewline(bs) {
				data = append(data, filterSymbols(line, b.escape)...)
				line = line[:0]
			}
		}
	} // end for
}
