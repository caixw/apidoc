// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"

	"github.com/issue9/assert"
)

var _ blocker = &swiftNestMCommentBlock{}

func TestSwiftNestCommentBlock(t *testing.T) {
	a := assert.New(t)

	b := newSwiftNestMCommentBlock("/*", "*/")
	a.NotNil(b)

	l := &lexer{data: []byte(`/* *123*123**/`)}
	a.True(b.BeginFunc(l))
	ret, ok := b.EndFunc(l)
	a.True(ok).
		Equal(string(ret), " *123*123*"). // 返回内容
		True(l.atEOF())                   // 到达末尾

	// 嵌套注释
	l = &lexer{data: []byte(`/*0/*1/*2*/*/*/`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(ret), "0/*1/*2*/*/"). // 返回内容
		True(l.atEOF())                    // 到达末尾

		// 多出 end 匹配项
	l = &lexer{data: []byte(`/*0/*1/*2*/*/*/*/`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(ret), "0/*1/*2*/*/"). // 返回内容
		Equal(string(l.data[l.pos:]), "*/")

	// 缺少 end 匹配项
	l = &lexer{data: []byte(`/*0/*1/*2*/*/`)}
	a.True(b.BeginFunc(l))
	ret, ok = b.EndFunc(l)
	a.False(ok).
		Equal(len(ret), 0).
		True(l.atEOF()) // 到达末尾
}
