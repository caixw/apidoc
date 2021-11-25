// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

func TestRubyMultipleComment(t *testing.T) {
	a := assert.New(t, false)
	b := newRubyMultipleComment("=pod", "=cut", "")

	rslt := messagetest.NewMessageHandler()
	l := newParser(rslt.Handler, core.Block{Data: []byte("=pod\ncomment1\n=cut\n")}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, found := b.endFunc(l)
	a.True(found).
		Equal(string(data), "     comment1\n     ")

	// 多个注释结束符
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte("=pod\ncomment1\ncomment2\n=cut\n=cut\n")}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, found = b.endFunc(l)
	a.True(found).
		Equal(string(data), "     comment1\ncomment2\n     ")

	// 换行符开头
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte("\ncomment1\ncomment2\n=cut\n=cut\n")}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.False(b.beginFunc(l))
	data, found = b.endFunc(l)
	a.True(found).
		Equal(string(data), "     \ncomment1\ncomment2\n     ")

	// 没有注释结束符
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte("comment1")}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.False(b.beginFunc(l))
	data, found = b.endFunc(l)
	a.False(found).Nil(data)
}
