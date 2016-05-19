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

	l := New([]rune(" line1\n @delimiter line2 \n"))
	a.NotNil(l)

	a.Equal(l.Read("@delimiter"), []rune("line1\n "))

	// 查找一个不存在的字符
	a.Equal(l.Read("not exists"), []rune("@delimiter line2 \n"))

	a.Equal(l.Read("end"), []rune(""))
}

func TestLexer_ReadWord(t *testing.T) {
	a := assert.New(t)

	l := New([]rune(" line1\n line2 \n"))
	a.NotNil(l)

	a.Equal(l.ReadWord(), []rune("line1"))
	a.Equal(l.ReadWord(), []rune("line2"))
}

func TestLexer_ReadLine(t *testing.T) {
	a := assert.New(t)

	l := New([]rune(" line1\n line2 \n"))
	a.NotNil(l)

	a.Equal(l.ReadLine(), []rune("line1"))
	a.Equal(l.ReadLine(), []rune("line2 "))
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

// go1.6 BenchmarkLexer_ReadWord-4	50000000	        34.5 ns/op
func BenchmarkLexer_ReadWord(b *testing.B) {
	a := assert.New(b)
	l := New([]rune("line1\n @delimiter line2 \n"))
	a.NotNil(l)

	for i := 0; i < b.N; i++ {
		_ = l.ReadWord()
		l.pos = 0
	}
}

// go1.6 BenchmarkLexer_ReadLine-4	100000000	        20.9 ns/op
func BenchmarkLexer_ReadLine(b *testing.B) {
	a := assert.New(b)
	l := New([]rune("line1\n @delimiter line2 \n"))
	a.NotNil(l)

	for i := 0; i < b.N; i++ {
		_ = l.ReadLine()
		l.pos = 0
	}
}

// go1.6 BenchmarkNew-4       	300000000	         5.66 ns/op
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New([]rune("line"))
	}
}
