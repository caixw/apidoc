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
//
//	/*
//	 *
//	 */
//
// 中的 * 字符
func newSwiftNestMCommentBlock(begin, end, prefix string) blocker {
	return &swiftNestMCommentBlock{
		begin:  begin,
		end:    end,
		prefix: []byte(prefix),
		begins: []byte(begin),
		ends:   []byte(end),
	}
}

func (b *swiftNestMCommentBlock) beginFunc(l *parser) bool {
	if l.Match(b.begin) {
		b.level++
		return true
	}

	return false
}

func (b *swiftNestMCommentBlock) endFunc(l *parser) (data []byte, ok bool) {
	data = append(make([]byte, 0, 200), b.begins...)

LOOP:
	for {
		switch {
		case l.AtEOF():
			return nil, false
		case l.Match(b.end):
			data = append(data, b.ends...)
			b.level--
			if b.level == 0 {
				break LOOP
			}
		case l.Match(b.begin):
			data = append(data, b.begins...)
			b.level++
		default:
			data = append(data, l.Next(1)...)
		}
	}

	return convertMultipleCommentToXML(data, b.begins, b.ends, b.prefix), true
}
