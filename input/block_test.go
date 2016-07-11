// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"

	"github.com/issue9/assert"
)

var _ blocker = &block{}

func TestBlock_BeginFunc_EndFunc(t *testing.T) {
	a := assert.New(t)
	bStr := &block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: "\\"}
	bSComment := &block{Type: blockTypeSComment, Begin: "//"}
	bMComment := &block{Type: blockTypeMComment, Begin: "/*", End: "*/"}

	l := &lexer{
		data: []byte("// scomment1\n// scomment2"),
	}
	a.False(bStr.BeginFunc(l))
	a.True(bSComment.BeginFunc(l))
	a.False(bMComment.BeginFunc(l))
	ret, ok := bSComment.EndFunc(l)
	a.True(ok).Equal(ret, []rune(" scomment1\n scomment2"))
	ret, ok = bMComment.EndFunc(l)
	a.False(ok).Equal(len(ret), 0)
}

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

func TestFilterSymbols(t *testing.T) {
	a := assert.New(t)

	eq := func(charset, v1, v2 string) {
		s1 := string(filterSymbols([]rune(v1), charset))
		a.Equal(s1, v2)
	}

	neq := func(charset, v1, v2 string) {
		s1 := string(filterSymbols([]rune(v1), charset))
		a.NotEqual(s1, v2)
	}

	eq("/*", "* line", "line")
	eq("/*", "*   line", "  line")
	eq("/*", "*\tline", "line")
	eq("/*", "* \tline", "\tline")
	eq("/*", "*\nline", "line")

	eq("/*", "/ line", "line")
	eq("/*", "/   line", "  line")

	eq("/*", "  * line", "line")
	eq("/*", "  *  line", " line")
	eq("/*", "\t*  line", " line")
	eq("/*", "\t* \nline", "\nline")
	eq("/*", "\t*\n line", " line")

	// 包含多个符号
	neq("/*", "// line", "line")
	neq("/*", "**   line", "  line")
	neq("/*", "/* line", "line")
	neq("/*", "*/   line", "  line")

	// 非定义的符号
	neq("/*", "+ line", "line")
	neq("/*", "+   line", "  line")

	eq("++", "+ line", "line")
	neq("++", "++ line", "line")
}
