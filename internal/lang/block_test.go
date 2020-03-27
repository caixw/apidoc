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

	l, err := NewLexer([]byte(`"text"`), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, ok := b.EndFunc(l)
	a.True(ok).Nil(data).Nil(raw)

	// 带转义字符
	l, err = NewLexer([]byte(`"te\"xt"`), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Nil(data).
		Nil(raw).
		Equal(l.Position().Offset, len(`"te\"xt"`))

	// 找不到匹配字符串
	l, err = NewLexer([]byte("text"), nil)
	a.NotError(err).NotNil(l)
	a.False(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.False(ok).Nil(data).Nil(raw)
}

func TestSingleComment(t *testing.T) {
	a := assert.New(t)
	b := newCStyleSingleComment()
	a.NotNil(b)

	l, err := NewLexer([]byte("//comment1\n"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, ok := b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "  comment1\n").
		Equal(string(raw), "//comment1\n")

	// 没有换行符，则自动取到结束符。
	l, err = NewLexer([]byte("// comment1"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "   comment1").
		Equal(string(raw), "// comment1")

	// 多行连续的单行注释，且 // 前带空格。
	l, err = NewLexer([]byte("//comment1\n//comment2\n // comment3"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "  comment1\n  comment2\n    comment3").
		Equal(string(raw), "//comment1\n//comment2\n // comment3")

	// 多行连续的单行注释，中间有空白行。
	l, err = NewLexer([]byte("//comment1\n//\n//comment2\n //comment3"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "  comment1\n  \n  comment2\n   comment3").
		Equal(string(raw), "//comment1\n//\n//comment2\n //comment3")

	// 多行不连续的单行注释。
	l, err = NewLexer([]byte("//comment1\n // comment2\n\n //comment3\n"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(data), "  comment1\n    comment2\n").
		Equal(string(raw), "//comment1\n // comment2\n")
}

func TestMultipleComment(t *testing.T) {
	a := assert.New(t)
	b := newCStyleMultipleComment()

	l, err := NewLexer([]byte("/*comment1\n*/"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, found := b.EndFunc(l)
	a.True(found).
		Equal(string(data), "  comment1\n  ").
		Equal(string(raw), "/*comment1\n*/")

	// 多个注释结束符
	l, err = NewLexer([]byte("/*comment1\ncomment2*/*/"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.True(found).
		Equal(string(data), "  comment1\ncomment2  ").
		Equal(string(raw), "/*comment1\ncomment2*/")

	// 空格开头
	l, err = NewLexer([]byte("/*\ncomment1\ncomment2*/*/"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.True(found).
		Equal(string(data), "  \ncomment1\ncomment2  ").
		Equal(string(raw), "/*\ncomment1\ncomment2*/")

	// 没有注释结束符
	l, err = NewLexer([]byte("comment1"), nil)
	a.NotError(err).NotNil(l)
	a.False(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.False(found).Nil(data).Nil(raw)
}

func TestConvertSingleCommentToXML(t *testing.T) {
	a := assert.New(t)
	data := []struct {
		input, begin, output string
	}{
		{},
		{
			input:  "// xxx",
			begin:  "/",
			output: "   xxx",
		},
		{
			input:  "//xxx",
			begin:  "//",
			output: "  xxx",
		},
		{
			input:  "\t//\txxx",
			begin:  "//",
			output: "\t  \txxx",
		},
		{
			input:  "#xxx",
			begin:  "#",
			output: " xxx",
		},
		{
			input:  "# xxx",
			begin:  "#",
			output: "  xxx",
		},
		{
			input:  "## xxx",
			begin:  "#",
			output: "   xxx",
		},
	}

	for i, item := range data {
		output := convertSingleCommentToXML([]byte(item.input), []byte(item.begin))
		a.Equal(string(output), item.output, "not equal @ %d,v1=%s,v2=%s", i, string(output), item.output)
	}
}

func TestConvertMultipleCommentToXML(t *testing.T) {
	a := assert.New(t)
	data := []struct {
		input, begin, end, chars, output string
	}{
		{},
		{
			input:  "/*\n * xx\n * xx\n */",
			begin:  "/*",
			end:    "*/",
			chars:  "*",
			output: "  \n   xx\n   xx\n   ",
		},
		{
			input:  "/**\n * xx\n * xx\n */",
			begin:  "/*",
			end:    "*/",
			chars:  "*",
			output: "   \n   xx\n   xx\n   ",
		},
		{
			input:  "/**xxx\n * xx\n * xx\n */",
			begin:  "/*",
			end:    "*/",
			chars:  "*",
			output: "  *xxx\n   xx\n   xx\n   ",
		},
		{
			input:  "/**xxx\n * xx\n * xx\n */",
			begin:  "/**",
			end:    "*/",
			chars:  "*",
			output: "   xxx\n   xx\n   xx\n   ",
		},
	}

	for i, item := range data {
		output := convertMultipleCommentToXML([]byte(item.input), []byte(item.begin), []byte(item.end), []byte(item.chars))
		a.Equal(string(output), item.output, "not equal @ %d,v1=%s,v2=%s", i, string(output), item.output)
	}
}

func TestReplaceSymbol(t *testing.T) {
	a := assert.New(t)
	data := []struct {
		input, chars, output string
	}{
		{},
		{
			input:  "// xxx",
			chars:  "/",
			output: "   xxx",
		},
		{
			input:  "/* xxx",
			chars:  "/*",
			output: "   xxx",
		},
		{
			input:  "/*xxx",
			chars:  "/*",
			output: "/*xxx",
		},
		{
			input:  " /* xxx",
			chars:  "/*",
			output: "    xxx",
		},
		{
			input:  "\t/*\txxx",
			chars:  "/*",
			output: "\t  \txxx",
		},
		{
			input:  "\t/**\txxx",
			chars:  "/*",
			output: "\t   \txxx",
		},
		{
			input:  "\t/**\n\txxx",
			chars:  "/*",
			output: "\t   \n\txxx",
		},
		{
			input:  "\t/**\n\t** xxx",
			chars:  "/*",
			output: "\t   \n\t   xxx",
		},
		{
			input:  "\t/**\n\t**xxx",
			chars:  "/*",
			output: "\t   \n\t**xxx",
		},
		{
			input:  "\t/**\n\t** xxx",
			chars:  "/*",
			output: "\t   \n\t   xxx",
		},
	}

	for i, item := range data {
		output := replaceSymbols([]byte(item.input), []byte(item.chars))
		a.Equal(string(output), item.output, "not equal @ %d,v1=%s,v2=%s", i, string(output), item.output)
	}
}
