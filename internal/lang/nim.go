// SPDX-License-Identifier: MIT

package lang

import "bytes"

type nimRawString struct {
	escape, begin1, begin2, end string
}

type nimMultipleString struct{}

func newNimRawString() blocker {
	return &nimRawString{
		escape: `""`,
		begin1: `r"`,
		begin2: `R"`,
		end:    `"`,
	}
}

func (s *nimRawString) beginFunc(l *parser) bool {
	return l.Match(s.begin1) || l.Match(s.begin2)
}

func (s *nimRawString) endFunc(l *parser) (data []byte, ok bool) {
	for {
		switch {
		case l.AtEOF():
			return nil, false
		case l.Match(s.escape): // 转义
			break
		case l.Match(s.end): // 结束
			return nil, true
		default:
			l.Next(1)
		}
	} // end for
}

func newNimMultipleString() blocker {
	return &nimMultipleString{}
}

func (s *nimMultipleString) beginFunc(l *parser) bool {
	return l.Match(`"""`)
}

func (s *nimMultipleString) endFunc(l *parser) ([]byte, bool) {
	for {
		pos := l.Current()
		switch {
		case l.AtEOF():
			return nil, false
		case l.Match(`"""`):
			if l.AtEOF() { // 后面没有内容了
				return nil, true
			}

			bs, ok := l.Delim('\n', true)
			if ok && len(bytes.TrimSpace(bs)) == 0 { // """ 后只有空格和回车符
				return nil, true
			}
			l.Move(pos)
			l.Next(1)
		default:
			l.Next(1)
		}
	}
}
