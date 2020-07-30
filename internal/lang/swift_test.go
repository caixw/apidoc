// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

var _ blocker = &swiftNestMCommentBlock{}

func TestSwiftNestCommentBlock(t *testing.T) {
	a := assert.New(t)

	b := newSwiftNestMCommentBlock("/*", "*/", "*")
	a.NotNil(b)

	rslt := messagetest.NewMessageHandler()
	l := newParser(rslt.Handler, core.Block{Data: []byte(`/* *123*123**/`)}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)

	a.True(b.beginFunc(l))
	data, ok := b.endFunc(l)
	a.True(ok).
		Equal(string(data), "   *123*123*  ")
	bs := l.Next(1)
	a.Empty(bs).True(l.AtEOF()) // 到达末尾

	// 多行，最后一行没有任何内容，则不返回数据
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`/**
	* xx
	* yy
*/`)}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.True(ok).
		Equal(string(data), "   \n\t  xx\n\t  yy\n  ")

	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`/**
	* xx/yy/zz
	* yy/zz/
	*/`)}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.True(ok).
		Equal(string(data), "   \n\t  xx/yy/zz\n\t  yy/zz/\n\t  ")

	// 嵌套注释
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`/*0/*1/*2*/*/*/`)}, nil)
	rslt.Handler.Stop()
	a.NotError(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.True(ok).
		Equal(string(data), "  0/*1/*2*/*/  ")
	bs = l.Next(1)
	a.Empty(bs).True(l.AtEOF()) // 到达末尾

	// 多出 end 匹配项
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`/*0/*1/*2*/*/*/*/`)}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.True(ok).
		Equal(string(data), "  0/*1/*2*/*/  ").
		Equal(string(l.All()), "*/")

	// 缺少 end 匹配项
	rslt = messagetest.NewMessageHandler()
	l = newParser(rslt.Handler, core.Block{Data: []byte(`/*0/*1/*2*/*/`)}, nil)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).NotNil(l)
	a.True(b.beginFunc(l))
	data, ok = b.endFunc(l)
	a.False(ok).
		Equal(len(data), 0).
		True(l.AtEOF()) // 到达末尾
}
