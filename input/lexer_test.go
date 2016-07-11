// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"
	"unicode/utf8"

	"github.com/issue9/assert"
)

func TestLexer_lineNumber(t *testing.T) {
	a := assert.New(t)

	l := &lexer{data: []byte("l0\nl1\nl2\nl3\n")}
	l.pos = 3
	a.Equal(l.lineNumber(), 1)

	l.pos += 3
	a.Equal(l.lineNumber(), 2)

	l.pos += 3
	l.pos += 3
	a.Equal(l.lineNumber(), 4)
}

func TestLexer_next(t *testing.T) {
	a := assert.New(t)

	l := &lexer{
		data: []byte("ab\ncd"),
	}

	a.Equal('a', l.next())
	a.Equal('b', l.next())
	a.Equal('\n', l.next())
	a.Equal('c', l.next())
	a.False(l.atEOF())
	a.Equal('d', l.next())
	a.True(l.atEOF())
	a.Equal(utf8.RuneError, l.next())
	a.Equal(utf8.RuneError, l.next())
	a.True(l.atEOF())
}

func TestLexer_match(t *testing.T) {
	a := assert.New(t)

	l := &lexer{
		data: []byte("ab\ncd"),
	}

	a.False(l.match("b")).Equal(0, l.pos)
	a.True(l.match("ab")).Equal(2, l.pos)

	l.pos = len(l.data)
	a.False(l.match("ab"))

	// 匹配结尾单词
	l.pos = 3 // c的位置
	a.True(l.match("cd"))
}

func TestLexer_block(t *testing.T) {
	a := assert.New(t)

	l := &lexer{
		data: []byte(`// scomment1
// scomment2
func(){}
"/*string1"
"//string2"
/*
mcomment1
mcomment2
*/

// scomment3
// scomment4
=pod
 mcomment3
 mcomment4
=cut
`),
	}

	blocks := []*block{
		&block{Type: blockTypeSComment, Begin: "//"},
		&block{Type: blockTypeMComment, Begin: "/*", End: "*/"},
		&block{Type: blockTypeMComment, Begin: "\n=pod", End: "\n=cut"},
		&block{Type: blockTypeString, Begin: `"`, End: `"`, Escape: "\\"},
	}

	b := l.block(blocks) // scomment1
	a.Equal(b.Type, blockTypeSComment)
	rs, err := b.end(l)
	a.NotError(err).Equal(string(rs), " scomment1\n scomment2\n")

	b = l.block(blocks) // string1
	a.Equal(b.Type, blockTypeString)
	_, err = b.end(l)
	a.NotError(err)

	b = l.block(blocks) // string2
	a.Equal(b.Type, blockTypeString)
	_, err = b.end(l)
	a.NotError(err)

	b = l.block(blocks)
	a.Equal(b.Type, blockTypeMComment) // mcomment1
	rs, err = b.end(l)
	a.NotError(err).Equal(string(rs), "\nmcomment1\nmcomment2\n")

	/* 测试一段单行注释后紧跟 \n=pod 形式的多行注释，是否会出错 */

	b = l.block(blocks) // scomment3,scomment4
	a.Equal(b.Type, blockTypeSComment)
	rs, err = b.end(l)
	a.NotError(err).Equal(string(rs), " scomment3\n scomment4\n")

	b = l.block(blocks) // mcomment3,mcomment4
	a.Equal(b.Type, blockTypeMComment)
	rs, err = b.end(l)
	a.NotError(err).Equal(string(rs), "\n mcomment3\n mcomment4")
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
