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

	l, err := NewLexer([]byte("=pod\ncomment1\n=cut\n"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, found := b.EndFunc(l)
	a.True(found).
		Equal(string(data), "     comment1\n     ").
		Equal(string(raw), "=pod\ncomment1\n=cut\n")

	// 多个注释结束符
	l, err = NewLexer([]byte("=pod\ncomment1\ncomment2\n=cut\n=cut\n"), nil)
	a.NotError(err).NotNil(l)
	a.True(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.True(found).
		Equal(string(data), "     comment1\ncomment2\n     ").
		Equal(string(raw), "=pod\ncomment1\ncomment2\n=cut\n")

	// 换行符开头
	l, err = NewLexer([]byte("\ncomment1\ncomment2\n=cut\n=cut\n"), nil)
	a.NotError(err).NotNil(l)
	a.False(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.True(found).
		Equal(string(data), "     \ncomment1\ncomment2\n     ").
		Equal(string(raw), "=pod\n\ncomment1\ncomment2\n=cut\n")

	// 没有注释结束符
	l, err = NewLexer([]byte("comment1"), nil)
	a.NotError(err).NotNil(l)
	a.False(b.BeginFunc(l))
	raw, data, found = b.EndFunc(l)
	a.False(found).Nil(data).Nil(raw)
}
