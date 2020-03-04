// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"
)

func TestPHPDocBlock(t *testing.T) {
	a := assert.New(t)
	b := newPHPDocBlock()
	a.NotNil(b)

	// herodoc
	l := &Lexer{data: []byte(`<<<EOF
	xx
	xx
EOF
`)}
	a.True(b.BeginFunc(l))
	bb, ok := b.(*phpDocBlock)
	a.True(ok)
	a.Equal(bb.token1, "\nEOF\n").
		Equal(bb.token2, "\nEOF;\n").
		Equal(bb.doctype, phpHerodoc)
	raw, data, ok := b.EndFunc(l)
	a.True(ok).
		Nil(data).
		Nil(raw)

	// nowdoc
	l = &Lexer{data: []byte(`<<<'EOF'
	xx
	xx
EOF
`)}
	a.True(b.BeginFunc(l))
	bb, ok = b.(*phpDocBlock)
	a.True(ok)
	a.Equal(bb.token1, "\nEOF\n").
		Equal(bb.token2, "\nEOF;\n").
		Equal(bb.doctype, phpNowdoc)
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Nil(data).
		Nil(raw)

	// nowdoc 验证结尾带分号的结束符
	l = &Lexer{data: []byte(`<<<'EOF'
	xx
	xx
EOF;
`)}
	a.True(b.BeginFunc(l))
	bb, ok = b.(*phpDocBlock)
	a.True(ok)
	a.Equal(bb.token1, "\nEOF\n").
		Equal(bb.token2, "\nEOF;\n").
		Equal(bb.doctype, phpNowdoc)
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Nil(raw).
		Nil(data)

	// 开始符号错误
	l = &Lexer{data: []byte(`<<<
	xx
	xx
EOF;
`)}
	a.False(b.BeginFunc(l))

	// nowdoc 不存在结束符
	l = &Lexer{data: []byte(`<<<'EOF'
	xx
	xx
EO
`)}
	a.True(b.BeginFunc(l))
	bb, ok = b.(*phpDocBlock)
	a.True(ok)
	a.Equal(bb.token1, "\nEOF\n").
		Equal(bb.token2, "\nEOF;\n").
		Equal(bb.doctype, phpNowdoc)
	raw, data, ok = b.EndFunc(l)
	a.False(ok)
}
