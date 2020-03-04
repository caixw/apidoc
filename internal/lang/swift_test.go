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

	l := &Lexer{data: []byte(`/* *123*123**/`)}
	a.True(b.BeginFunc(l))
	ret, ok := b.EndFunc(l)
	a.True(ok).
		Equal(string(ret), " *123*123*"). // 返回内容
		True(l.AtEOF())                   // 到达末尾

	// 多行，最后一行没有任何内容，则不返回数据
	l = &Lexer{data: []byte(`/**
	* xx
	* yy
*/`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(ret), "\n xx\n yy\n")

	l = &Lexer{data: []byte(`/**
	* xx/yy/zz
	* yy/zz/
	*/`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(ret), "\n xx/yy/zz\n yy/zz/\n\t")

	// 嵌套注释
	l = &Lexer{data: []byte(`/*0/*1/*2*/*/*/`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(ret), "0/*1/*2*/*/"). // 返回内容
		True(l.AtEOF())                    // 到达末尾

		// 多出 end 匹配项
	l = &Lexer{data: []byte(`/*0/*1/*2*/*/*/*/`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(ret), "0/*1/*2*/*/"). // 返回内容
		Equal(string(l.data[l.pos:]), "*/")

	// 缺少 end 匹配项
	l = &Lexer{data: []byte(`/*0/*1/*2*/*/`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.False(ok).
		Equal(len(ret), 0).
		True(l.AtEOF()) // 到达末尾
}
