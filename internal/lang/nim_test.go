// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

func TestNimRawString(t *testing.T) {
	a := assert.New(t)

	b := newNimRawString()
	a.NotNil(b)

	rslt := messagetest.NewMessageHandler()
	l := newParser(rslt.Handler, core.Block{Data: []byte(`r"123""123"`)}, nil)
	a.NotNil(l)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok := b.endFunc(l)
	a.True(ok).
		Equal(len(data), 0) // 不返回内容
	bs := l.Next(1)             // 继续向后推进，才会
	a.Empty(bs).True(l.AtEOF()) // 到达末尾

	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`R"123123"`)}, nil)
	a.NotNil(l)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.True(ok).
		Empty(data) // 不返回内容
	bs = l.Next(1)              // 继续向后推进，才会
	a.Empty(bs).True(l.AtEOF()) // 到达末尾
}

func TestNimMultipleString(t *testing.T) {
	a := assert.New(t)

	b := newNimMultipleString()
	a.NotNil(b)

	rslt := messagetest.NewMessageHandler()
	l := newParser(rslt.Handler, core.Block{Data: []byte(`"""line1
	a.NotNil(l)
line2
"""`)}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok := b.endFunc(l)
	a.True(ok).Nil(data)

	// 结束符后带空格
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`"""line1
line2
"""  
`)}, nil)
	a.NotNil(l)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.True(ok).Nil(data)

	// 结束符后带字符
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`"""line1
line2
""" suffix 
`)}, nil)
	a.NotNil(l)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.False(ok).Nil(data)
}
