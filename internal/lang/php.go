// SPDX-License-Identifier: MIT

package lang

const (
	phpHerodoc int8 = iota + 1
	phpNowdoc
)

type phpDocBlock struct {
	token1  string
	token2  string
	doctype int8
}

// herodoc 和 nowdoc 的实现。
//
// http://php.net/manual/zh/language.types.string.php#language.types.string.syntax.heredoc
func newPHPDocBlock() Blocker {
	return &phpDocBlock{
		doctype: phpHerodoc,
	}
}

func (b *phpDocBlock) BeginFunc(l *Lexer) bool {
	if !l.match("<<<") {
		return false
	}

	token := readLine(l)
	if len(token) == 0 {
		l.pos -= 3 // 退回 <<< 字符
		return false
	}

	if token[0] == '\'' && token[len(token)-1] == '\'' {
		b.doctype = phpNowdoc
		token = token[1 : len(token)-1]
	}

	b.token1 = "\n" + string(token) + "\n"
	b.token2 = "\n" + string(token) + ";\n"

	return true
}

func (b *phpDocBlock) EndFunc(l *Lexer) ([][]byte, bool) {
	for {
		switch {
		case l.AtEOF():
			return nil, false
		case l.match(b.token1):
			return nil, true
		case l.match(b.token2):
			return nil, true
		default:
			l.pos++
		}
	}
}

// 读取到当前行行尾。
//
// 返回 nil 表示没有换行符，即当前就是最后一行。
func readLine(l *Lexer) []byte {
	start := l.pos
	for index, b := range l.data[l.pos:] {
		if l.AtEOF() {
			return nil
		}

		if b == '\n' {
			l.pos += index
			return l.data[start : index+start]
		}
	} // end for

	return nil
}
