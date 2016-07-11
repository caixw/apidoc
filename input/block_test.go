// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"

	"github.com/issue9/assert"
)

func TestBlock_endString(t *testing.T) {
	a := assert.New(t)
	b := &block{
		Type:   blockTypeString,
		Begin:  `"`,
		End:    `"`,
		Escape: "\\",
	}

	l := &lexer{
		data: []byte(`text"`),
	}
	rs, ok := b.endString(l)
	a.True(ok).Nil(rs)

	// 带转义字符
	l = &lexer{
		data: []byte(`te\"xt"`),
	}
	rs, ok = b.endString(l)
	a.True(ok).
		Nil(rs).
		Equal(l.pos, len(l.data))

	// 找不到匹配字符串
	l = &lexer{
		data: []byte("text"),
	}
	rs, ok = b.endString(l)
	a.False(ok).Nil(rs)
}

func TestBlock_endSComment(t *testing.T) {
	a := assert.New(t)
	b := &block{
		Type:  blockTypeSComment,
		Begin: `//`,
	}

	l := &lexer{
		data: []byte("comment1\n"),
	}
	rs, err := b.endSComments(l)
	a.NotError(err).Equal(string(rs), "comment1\n")

	// 没有换行符，则自动取到结束符。
	l = &lexer{
		data: []byte("comment1"),
	}
	rs, err = b.endSComments(l)
	a.NotError(err).Equal(string(rs), "comment1")

	// 多行连续的单行注释。
	l = &lexer{
		data: []byte("comment1\n//comment2\n //comment3"),
	}
	rs, err = b.endSComments(l)
	a.NotError(err).Equal(string(rs), "comment1\ncomment2\ncomment3")

	// 多行不连续的单行注释。
	l = &lexer{
		data: []byte("comment1\n // comment2\n\n //comment3\n"),
	}
	rs, err = b.endSComments(l)
	a.NotError(err).Equal(string(rs), "comment1\n comment2\n")
}

func TestBlock_endMComment(t *testing.T) {
	a := assert.New(t)
	b := &block{
		Type:  blockTypeSComment,
		Begin: "/*",
		End:   "*/",
	}

	l := &lexer{
		data: []byte("comment1\n*/"),
	}
	rs, found := b.endMComments(l)
	a.True(found).Equal(string(rs), "comment1\n")

	// 多个注释结束符
	l = &lexer{
		data: []byte("comment1\ncomment2*/*/"),
	}
	rs, found = b.endMComments(l)
	a.True(found).Equal(string(rs), "comment1\ncomment2")

	// 空格开头
	l = &lexer{
		data: []byte("\ncomment1\ncomment2*/*/"),
	}
	rs, found = b.endMComments(l)
	a.True(found).Equal(string(rs), "\ncomment1\ncomment2")

	// 没有注释结束符
	l = &lexer{
		data: []byte("comment1"),
	}
	rs, found = b.endMComments(l)
	a.False(found).Nil(rs)
}

func TestBlock_filterSymbols(t *testing.T) {
	a := assert.New(t)
	b := &block{Begin: "/*"}

	eq := func(b *block, v1, v2 string) {
		s1 := string(b.filterSymbols([]rune(v1)))
		a.Equal(s1, v2)
	}

	neq := func(b *block, v1, v2 string) {
		s1 := string(b.filterSymbols([]rune(v1)))
		a.NotEqual(s1, v2)
	}

	eq(b, "* line", "line")
	eq(b, "*   line", "  line")
	eq(b, "*\tline", "line")
	eq(b, "* \tline", "\tline")
	eq(b, "*\nline", "line")

	eq(b, "/ line", "line")
	eq(b, "/   line", "  line")

	eq(b, "  * line", "line")
	eq(b, "  *  line", " line")
	eq(b, "\t*  line", " line")
	eq(b, "\t* \nline", "\nline")
	eq(b, "\t*\n line", " line")

	// 包含多个符号
	neq(b, "// line", "line")
	neq(b, "**   line", "  line")
	neq(b, "/* line", "line")
	neq(b, "*/   line", "  line")

	// 非定义的符号
	neq(b, "+ line", "line")
	neq(b, "+   line", "  line")

	b = &block{Begin: "++"}
	eq(b, "+ line", "line")
	neq(b, "++ line", "line")
}
