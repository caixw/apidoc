// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

func TestPHPDocBlock(t *testing.T) {
	a := assert.New(t, false)
	b := newPHPDocBlock()
	a.NotNil(b)

	// herodoc
	data := []byte(`<<<EOF
	xx
	xx
EOF
`)

	rslt := messagetest.NewMessageHandler()
	l := newParser(rslt.Handler, core.Block{Data: data}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	bb, ok := b.(*phpDocBlock)
	a.True(ok)
	a.Equal(bb.token1, "\nEOF\n").
		Equal(bb.token2, "\nEOF;\n").
		Equal(bb.doctype, phpHerodoc)
	data, ok = b.endFunc(l)
	a.True(ok).
		Nil(data)

	// nowdoc
	data = []byte(`<<<'EOF'
	xx
	xx
EOF
`)
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: data}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	bb, ok = b.(*phpDocBlock)
	a.True(ok)
	a.Equal(bb.token1, "\nEOF\n").
		Equal(bb.token2, "\nEOF;\n").
		Equal(bb.doctype, phpNowdoc)
	data, ok = b.endFunc(l)
	a.True(ok).
		Nil(data)

	// nowdoc 验证结尾带分号的结束符
	data = []byte(`<<<'EOF'
	xx
	xx
EOF;
`)
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: data}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	bb, ok = b.(*phpDocBlock)
	a.True(ok)
	a.Equal(bb.token1, "\nEOF\n").
		Equal(bb.token2, "\nEOF;\n").
		Equal(bb.doctype, phpNowdoc)
	data, ok = b.endFunc(l)
	a.True(ok).
		Nil(data)

	// 开始符号错误
	data = []byte(`<<<
	xx
	xx
EOF;
`)
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: data}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.False(b.beginFunc(l))

	// nowdoc 不存在结束符
	data = []byte(`<<<'EOF'
	xx
	xx
EO
`)
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: data}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	bb, ok = b.(*phpDocBlock)
	a.True(ok)
	a.Equal(bb.token1, "\nEOF\n").
		Equal(bb.token2, "\nEOF;\n").
		Equal(bb.doctype, phpNowdoc)
	data, ok = b.endFunc(l)
	a.False(ok)
}
