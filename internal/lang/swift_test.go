// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

var _ Blocker = &swiftNestMCommentBlock{}

func TestSwiftNestCommentBlock(t *testing.T) {
	a := assert.New(t)

	b := newSwiftNestMCommentBlock("/*", "*/", "*")
	a.NotNil(b)

	l, err := NewLexer(core.Block{Data: []byte(`/* *123*123**/`)}, nil)
	a.NotError(err).NotNil(l)

	a.True(b.BeginFunc(l))
	data, ok := b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "   *123*123*  ")
	bs := l.Next(1)
	a.Empty(bs).True(l.AtEOF()) // 到达末尾

	// 多行，最后一行没有任何内容，则不返回数据
	l, err = NewLexer(core.Block{Data: []byte(`/**
	* xx
	* yy
*/`)}, nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "   \n\t  xx\n\t  yy\n  ")

	l, err = NewLexer(core.Block{Data: []byte(`/**
	* xx/yy/zz
	* yy/zz/
	*/`)}, nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "   \n\t  xx/yy/zz\n\t  yy/zz/\n\t  ")

	// 嵌套注释
	l, err = NewLexer(core.Block{Data: []byte(`/*0/*1/*2*/*/*/`)}, nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "  0/*1/*2*/*/  ")
	bs = l.Next(1)
	a.Empty(bs).True(l.AtEOF()) // 到达末尾

	// 多出 end 匹配项
	l, err = NewLexer(core.Block{Data: []byte(`/*0/*1/*2*/*/*/*/`)}, nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "  0/*1/*2*/*/  ").
		Equal(string(l.All()), "*/")

	// 缺少 end 匹配项
	l, err = NewLexer(core.Block{Data: []byte(`/*0/*1/*2*/*/`)}, nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	data, ok = b.EndFunc(l)
	a.False(ok).
		Equal(len(data), 0).
		True(l.AtEOF()) // 到达末尾
}
