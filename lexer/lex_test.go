// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLexer_lineNumber(t *testing.T) {
	a := assert.New(t)
	l := New([]rune("\n\n"))
	a.NotNil(l)

	a.Equal(0, l.lineNumber())

	l.pos++
	a.Equal(1, l.lineNumber())

	l.pos++
	a.Equal(2, l.lineNumber())
}

func TestLexer_Match(t *testing.T) {
	a := assert.New(t)
	l := New([]rune("line1\n line2 \n"))
	a.NotNil(l)

	a.True(l.Match("line"))
	a.Equal('1', l.data[l.pos])
	l.pos++

	l.pos++ // \n
	l.pos++ // 空格
	l.pos++ // l

	a.False(l.Match("2222")) // 不匹配，不会移动位置
	a.True(l.Match("ine2"))  // 正确匹配
	l.Backup()
	l.Backup()
	a.Equal('i', l.data[l.pos])
	l.pos++

	// 超过剩余字符的长度。
	a.False(l.Match("ne2\n\n"))
}

func TestLexer_Read(t *testing.T) {
	a := assert.New(t)

	l := New([]rune("line1\n @delimiter line2 \n"))
	a.NotNil(l)

	word := l.Read("@delimiter")
	a.Equal(string(word), "line1\n ")

	// 查找一个不存在的字符
	word = l.Read("not exists")
	a.Equal(string(word), "@delimiter line2 \n")

	word = l.Read("end")
	a.Equal(string(word), "")
}

func TestLexer_SkipSpace(t *testing.T) {
	a := assert.New(t)

	l := New([]rune(" line1\n line2 \n"))
	a.NotNil(l)

	l.SkipSpace()
	a.Equal(l.data[l.pos:], "line1\n line2 \n")

	a.True(l.Match("line1"))
	a.Equal(l.data[l.pos:], "\n line2 \n")
	l.SkipSpace()
	a.Equal(l.data[l.pos:], "line2 \n")
}

// go1.6 BenchmarkLexer_Read-4	10000000	       130 ns/op
func BenchmarkLexer_Read(b *testing.B) {
	a := assert.New(b)
	l := New([]rune("line1\n @delimiter line2 \n"))
	a.NotNil(l)

	for i := 0; i < b.N; i++ {
		_ = l.Read("@delimiter")
		l.pos = 0
	}
}

// go1.6 BenchmarkNew-4       	300000000	         5.66 ns/op
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New([]rune("line"))
	}
}
