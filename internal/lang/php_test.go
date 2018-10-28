// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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
	l := &lexer{data: []byte(`<<<EOF
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
	ret, ok := b.EndFunc(l)
	a.True(ok).
		Nil(ret)

	// nowdoc
	l = &lexer{data: []byte(`<<<'EOF'
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
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Nil(ret)

	// nowdoc 验证结尾带分号的结束符
	l = &lexer{data: []byte(`<<<'EOF'
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
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Nil(ret)

	// 开始符号错误
	l = &lexer{data: []byte(`<<<
	xx
	xx
EOF;
`)}
	a.False(b.BeginFunc(l))

	// nowdoc 不存在结束符
	l = &lexer{data: []byte(`<<<'EOF'
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
	ret, ok = b.EndFunc(l)
	a.False(ok)
}

func TestReadLine(t *testing.T) {
	a := assert.New(t)

	l := &lexer{data: []byte("123")}
	a.Nil(readLine(l))

	l = &lexer{data: []byte("123\n")}
	a.Equal(string(readLine(l)), "123").
		Equal(l.pos, 3)

	l = &lexer{data: []byte("123\n"), pos: 1}
	a.Equal(string(readLine(l)), "23").
		Equal(l.pos, 3)
}
