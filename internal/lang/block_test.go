// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"
)

var _ Blocker = &block{}

func TestBlock_BeginFunc_EndFunc(t *testing.T) {
	a := assert.New(t)
	bStr := &block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: "\\"}
	bSComment := &block{Type: blockTypeSComment, Begin: "//"}
	bMComment := &block{Type: blockTypeMComment, Begin: "/*", End: "*/"}

	l := &Lexer{
		data: []byte("// scomment1\n// scomment2"),
	}
	a.False(bStr.BeginFunc(l))
	a.True(bSComment.BeginFunc(l))
	a.False(bMComment.BeginFunc(l))
	raw, data, ok := bSComment.EndFunc(l)
	a.True(ok).
		Equal(string(data), " scomment1\n scomment2").
		Equal(string(raw), " scomment1\n scomment2")
	raw, data, ok = bMComment.EndFunc(l)
	a.False(ok).Equal(len(data), 0)
}

func TestBlock_endString(t *testing.T) {
	a := assert.New(t)
	b := &block{
		Type:   blockTypeString,
		Begin:  `"`,
		End:    `"`,
		Escape: "\\",
	}

	l := &Lexer{
		data: []byte(`text"`),
	}
	raw, data, ok := b.endString(l)
	a.True(ok).Nil(data).Nil(raw)

	// 带转义字符
	l = &Lexer{
		data: []byte(`te\"xt"`),
	}
	raw, data, ok = b.endString(l)
	a.True(ok).
		Nil(data).
		Nil(raw).
		Equal(l.offset, len(l.data))

	// 找不到匹配字符串
	l = &Lexer{
		data: []byte("text"),
	}
	raw, data, ok = b.endString(l)
	a.False(ok).Nil(data).Nil(raw)
}

func TestBlock_endSComment(t *testing.T) {
	a := assert.New(t)
	b := &block{
		Type:  blockTypeSComment,
		Begin: `//`,
	}

	l := &Lexer{
		data: []byte("comment1\n"),
	}
	raw, data, err := b.endSComments(l)
	a.NotError(err).
		Equal(string(data), "comment1\n").
		Equal(string(raw), "comment1\n")

	// 没有换行符，则自动取到结束符。
	l = &Lexer{
		data: []byte("comment1"),
	}
	raw, data, err = b.endSComments(l)
	a.NotError(err).
		Equal(string(data), "comment1").
		Equal(string(raw), "comment1")

	// 多行连续的单行注释。
	l = &Lexer{
		data: []byte("comment1\n//comment2\n //comment3"),
	}
	raw, data, err = b.endSComments(l)
	a.NotError(err).
		Equal(string(data), "comment1\ncomment2\ncomment3").
		Equal(string(raw), "comment1\ncomment2\ncomment3")

	// 多行连续的单行注释，中间有空白行。
	l = &Lexer{
		data: []byte("comment1\n//\n//comment2\n //comment3"),
	}
	raw, data, err = b.endSComments(l)
	a.NotError(err).
		Equal(string(data), "comment1\n\ncomment2\ncomment3").
		Equal(string(raw), "comment1\n\ncomment2\ncomment3")

	// 多行不连续的单行注释。
	l = &Lexer{
		data: []byte("comment1\n // comment2\n\n //comment3\n"),
	}
	raw, data, err = b.endSComments(l)
	a.NotError(err).
		Equal(string(data), "comment1\n comment2\n").
		Equal(string(raw), "comment1\n comment2\n")
}

func TestBlock_endMComment(t *testing.T) {
	a := assert.New(t)
	b := &block{
		Type:  blockTypeSComment,
		Begin: "/*",
		End:   "*/",
	}

	l := &Lexer{
		data: []byte("comment1\n*/"),
	}
	raw, data, found := b.endMComments(l)
	a.True(found).
		Equal(string(data), "comment1\n").
		Equal(string(raw), "comment1\n")

	// 多个注释结束符
	l = &Lexer{
		data: []byte("comment1\ncomment2*/*/"),
	}
	raw, data, found = b.endMComments(l)
	a.True(found).
		Equal(string(data), "comment1\ncomment2").
		Equal(string(raw), "comment1\ncomment2")

	// 空格开头
	l = &Lexer{
		data: []byte("\ncomment1\ncomment2*/*/"),
	}
	raw, data, found = b.endMComments(l)
	a.True(found).
		Equal(string(data), "\ncomment1\ncomment2").
		Equal(string(raw), "\ncomment1\ncomment2")

	// 没有注释结束符
	l = &Lexer{
		data: []byte("comment1"),
	}
	raw, data, found = b.endMComments(l)
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
