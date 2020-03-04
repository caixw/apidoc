// SPDX-License-Identifier: MIT

package lang

// 描述了 pascal/delphi 语言的字符串，在 pascal 中
// 转义字符即引号本身，不适合直接在 block 中定义。
type pascalStringBlock struct {
	symbol string
	escape string
}

func newPascalStringBlock(symbol byte) Blocker {
	s := string(symbol)
	return &pascalStringBlock{
		symbol: s,
		escape: s + s,
	}
}

func (b *pascalStringBlock) BeginFunc(l *Lexer) bool {
	return l.match(b.symbol)
}

func (b *pascalStringBlock) EndFunc(l *Lexer) ([]byte, bool) {
LOOP:
	for {
		switch {
		case l.AtEOF():
			return nil, false
		case l.match(b.escape): // 转义
			continue LOOP
		case l.match(b.symbol): // 结束
			return nil, true
		default:
			l.pos++
		}
	} // end for
}
