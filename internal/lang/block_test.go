// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"
)

var (
	_ Blocker = &stringBlock{}
	_ Blocker = &singleComment{}
	_ Blocker = &multipleComment{}
)

func TestStringBlock(t *testing.T) {
	a := assert.New(t)
	b := newCStyleString()
	a.NotNil(b)

	l := &Lexer{
		data: []byte(`"text"`),
	}
	a.True(b.BeginFunc(l))
	raw, data, ok := b.EndFunc(l)
	a.True(ok).Nil(data).Nil(raw)

	// 带转义字符
	l = &Lexer{
		data: []byte(`"te\"xt"`),
	}
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Nil(data).
		Nil(raw).
		Equal(l.current.Offset, len(l.data))

	// 找不到匹配字符串
	l = &Lexer{
		data: []byte("text"),
	}
	a.False(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.False(ok).Nil(data).Nil(raw)
}

func TestSingleComment(t *testing.T) {
	a := assert.New(t)
	b := newCStyleSingleComment()
	a.NotNil(b)

	l := &Lexer{
		data: []byte("//comment1\n"),
	}
	a.True(b.BeginFunc(l))
	raw, data, err := b.EndFunc(l)
	a.NotError(err).
		Equal(string(data), "comment1\n").
		Equal(string(raw), "//comment1\n")

	// 没有换行符，则自动取到结束符。
	l = &Lexer{
		data: []byte("// comment1"),
	}
	a.True(b.BeginFunc(l))
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(data), " comment1").
		Equal(string(raw), "// comment1")

	// 多行连续的单行注释，且 // 前带空格。
	l = &Lexer{
		data: []byte("//comment1\n//comment2\n // comment3"),
	}
	a.True(b.BeginFunc(l))
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(data), "comment1\ncomment2\n comment3").
		Equal(string(raw), "//comment1\n//comment2\n // comment3")

	// 多行连续的单行注释，中间有空白行。
	l = &Lexer{
		data: []byte("//comment1\n//\n//comment2\n //comment3"),
	}
	a.True(b.BeginFunc(l))
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(data), "comment1\n\ncomment2\ncomment3").
		Equal(string(raw), "//comment1\n//\n//comment2\n //comment3")

	// 多行不连续的单行注释。
	l = &Lexer{
		data: []byte("//comment1\n // comment2\n\n //comment3\n"),
	}
	a.True(b.BeginFunc(l))
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(data), "comment1\n comment2\n").
		Equal(string(raw), "//comment1\n // comment2\n")
}

func TestMultipleComment(t *testing.T) {
	a := assert.New(t)
	b := newCStyleMultipleComment()

	l := &Lexer{
		data: []byte("/*comment1\n*/"),
	}
	a.True(b.BeginFunc(l))
	raw, data, found := b.EndFunc(l)
	a.True(found).
		Equal(string(data), "comment1\n").
		Equal(string(raw), "/*comment1\n*/")

	// 多个注释结束符
	l = &Lexer{
		data: []byte("/*comment1\ncomment2*/*/"),
	}
	a.True(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.True(found).
		Equal(string(data), "comment1\ncomment2").
		Equal(string(raw), "/*comment1\ncomment2*/")

	// 空格开头
	l = &Lexer{
		data: []byte("/*\ncomment1\ncomment2*/*/"),
	}
	a.True(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.True(found).
		Equal(string(data), "\ncomment1\ncomment2").
		Equal(string(raw), "/*\ncomment1\ncomment2*/")

	// 没有注释结束符
	l = &Lexer{
		data: []byte("comment1"),
	}
	a.False(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.False(found).Nil(data).Nil(raw)
}

func TestFilterSymbols(t *testing.T) {
	a := assert.New(t)

	eq := func(charset, v1, v2 string) {
		s1 := string(filterSymbols([]byte(v1), charset))
		a.Equal(s1, v2)
	}

	neq := func(charset, v1, v2 string) {
		s1 := string(filterSymbols([]byte(v1), charset))
		a.NotEqual(s1, v2)
	}

	eq("/*", "* ", " ")
	eq("/*", "* line", " line")
	eq("/*", "** line", " line")
	eq("/*", "*line", "*line")
	eq("/*", "**line", "**line")
	eq("/*", "* line", " line")
	eq("/*", " * line", " line")
	eq("/*", "*** line", " line")
	eq("/*", "*   line", "   line")
	eq("/*", "*\tline", "\tline")
	eq("/*", "* \tline", " \tline")

	eq("/*", "/ line", " line")
	eq("/*", "/   line", "   line")

	eq("/*", "  * line", " line")
	eq("/*", "  *  line", "  line")
	eq("/*", "\t*  line", "  line")

	neq("/*", "*\nline", "line")
	// 包含多个符号
	neq("/*", "// line", "line")
	neq("/*", "**   line", "  line")
	neq("/*", "/* line", "line")
	neq("/*", "*/   line", "  line")

	// 非定义的符号
	neq("/*", "+ line", "line")
	neq("/*", "+   line", "  line")

	eq("++", "+ line", " line")
	neq("++", "++ line", "line")
}
