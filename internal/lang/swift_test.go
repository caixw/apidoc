// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"
)

var _ Blocker = &swiftNestMCommentBlock{}

func TestSwiftNestCommentBlock(t *testing.T) {
	a := assert.New(t)

	b := newSwiftNestMCommentBlock("/*", "*/", "*")
	a.NotNil(b)

	l := NewLexer([]byte(`/* *123*123**/`), nil)
	a.True(b.BeginFunc(l))
	raw, data, ok := b.EndFunc(l)
	a.True(ok).
		Equal(string(data), " *123*123*"). // 返回内容
		Equal(string(raw), " *123*123*").  // 返回内容
		True(l.AtEOF())                    // 到达末尾

	// 多行，最后一行没有任何内容，则不返回数据
	l = NewLexer([]byte(`/**
	* xx
	* yy
*/`), nil)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "\n xx\n yy\n").
		Equal(string(raw), "*\n\t* xx\n\t* yy\n")

	l = NewLexer([]byte(`/**
	* xx/yy/zz
	* yy/zz/
	*/`), nil)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "\n xx/yy/zz\n yy/zz/\n\t").
		Equal(string(raw), "*\n\t* xx/yy/zz\n\t* yy/zz/\n\t")

	// 嵌套注释
	l = NewLexer([]byte(`/*0/*1/*2*/*/*/`), nil)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "0/*1/*2*/*/"). // 返回内容
		Equal(string(raw), "0/*1/*2*/*/").  // 返回内容
		True(l.AtEOF())                     // 到达末尾

		// 多出 end 匹配项
	l = NewLexer([]byte(`/*0/*1/*2*/*/*/*/`), nil)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "0/*1/*2*/*/"). // 返回内容
		Equal(string(raw), "0/*1/*2*/*/").  // 返回内容
		Equal(string(l.data[l.offset:]), "*/")

	// 缺少 end 匹配项
	l = NewLexer([]byte(`/*0/*1/*2*/*/`), nil)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.False(ok).
		Equal(len(data), 0).
		Equal(len(raw), 0).
		True(l.AtEOF()) // 到达末尾
}
