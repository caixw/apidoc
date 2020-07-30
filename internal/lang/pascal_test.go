// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

var _ blocker = &pascalStringBlock{}

func TestPascalStringBlock(t *testing.T) {
	a := assert.New(t)

	b := newPascalStringBlock('"')
	a.NotNil(b)

	rslt := messagetest.NewMessageHandler()
	l := newParser(rslt.Handler, core.Block{Data: []byte(`"123""123"`)}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok := b.endFunc(l)
	a.True(ok).
		Equal(len(data), 0) // 不返回内容
	bs := l.Next(1)             // 继续向后推进，才会
	a.Empty(bs).True(l.AtEOF()) // 到达末尾

	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`"123"""123"`)}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.True(ok).
		Equal(len(data), 0).            // 不返回内容
		Equal(string(l.All()), "123\"") // 未到达末尾
}
