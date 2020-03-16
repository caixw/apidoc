// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"
)

var _ Blocker = &rubyMultipleComment{}

func TestRubyMultipleComment(t *testing.T) {
	a := assert.New(t)
	b := newRubyMultipleComment("=pod", "=cut", "")

	l := &Lexer{
		data: []byte("=pod\ncomment1\n=cut\n"),
	}
	a.True(b.BeginFunc(l))
	raw, data, found := b.EndFunc(l)
	a.True(found).
		Equal(string(data), "     comment1\n     ").
		Equal(string(raw), "=pod\ncomment1\n=cut\n")

	// 多个注释结束符
	l = &Lexer{
		data: []byte("=pod\ncomment1\ncomment2\n=cut\n=cut\n"),
	}
	a.True(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.True(found).
		Equal(string(data), "     comment1\ncomment2\n     ").
		Equal(string(raw), "=pod\ncomment1\ncomment2\n=cut\n")

	// 换行符开头
	l = &Lexer{
		data: []byte("\ncomment1\ncomment2\n=cut\n=cut\n"),
	}
	a.False(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.True(found).
		Equal(string(data), "     \ncomment1\ncomment2\n     ").
		Equal(string(raw), "=pod\n\ncomment1\ncomment2\n=cut\n")

	// 没有注释结束符
	l = &Lexer{
		data: []byte("comment1"),
	}
	a.False(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.False(found).Nil(data).Nil(raw)
}
