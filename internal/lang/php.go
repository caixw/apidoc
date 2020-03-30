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
	prev := l.Position()

	if !l.Match("<<<") {
		return false
	}

	token, found := l.Delim('\n', true)
	if !found || len(token) <= 1 { // <<< 之后直接是换行符，则应该退回 <<< 字符
		l.SetPosition(prev)
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
		case l.AtEOF():
			return nil, nil, false
		case l.Match(b.token1):
			return nil, nil, true
		case l.Match(b.token2):
			return nil, nil, true
		default:
			l.Next(1)
		}
	}
}
