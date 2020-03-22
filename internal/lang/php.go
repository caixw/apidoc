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

	prev := l.prev
	token := l.delim('\n')
	if len(token) == 1 { // <<< 之后直接是换行符，则应该退回 <<< 字符
		l.current = prev
		return false
	}
	token = token[:len(token)-1] // l.delim 会带上换行符，需要去掉

	if token[0] == '\'' && token[len(token)-1] == '\'' {
		b.doctype = phpNowdoc
		token = token[1 : len(token)-1]
	}

	b.token1 = "\n" + string(token) + "\n"
	b.token2 = "\n" + string(token) + ";\n"

	return true
}

func (b *phpDocBlock) EndFunc(l *Lexer) (raw, data []byte, ok bool) {
	for {
		switch {
		case l.atEOF:
			return nil, nil, false
		case l.match(b.token1):
			return nil, nil, true
		case l.match(b.token2):
			return nil, nil, true
		default:
			l.next(1)
		}
	}
}
