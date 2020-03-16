// SPDX-License-Identifier: MIT

package lang

// 表示超始和结束符号必须占满一行的情况
type rubyMultipleComment struct {
	begin, end           string
	begins, ends, prefix []byte
}

func newRubyMultipleComment(begin, end, preifx string) Blocker {
	begin += "\n"
	end += "\n"
	return &rubyMultipleComment{
		begin:  begin,
		end:    end,
		prefix: []byte(preifx),
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
	raw = append(make([]byte, 0, 200), b.begins...)

LOOP:
	for {
		switch {
		case l.AtEOF(): // 没有找到结束符号，直接到达文件末尾
			return nil, nil, false
		case l.Position().Character == 0 && l.match(b.end):
			raw = append(raw, b.ends...)
			break LOOP
		default:
			raw = append(raw, l.next(1)...)
		}
	} // end for
	return raw, convertMultipleCommentToXML(raw, b.begins, b.ends, b.prefix), true
}
