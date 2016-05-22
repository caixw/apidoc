// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLexer_lineNumber(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]rune("\n\n"))
	a.NotNil(l)

	a.Equal(0, l.lineNumber())

	l.pos++
	a.Equal(1, l.lineNumber())

	l.pos++
	a.Equal(2, l.lineNumber())
}

func TestLexer_match(t *testing.T) {
	a := assert.New(t)
	l := newLexer([]rune("line1\n line2 \n"))
	a.NotNil(l)

	a.True(l.match("Line"))
	a.Equal('1', l.data[l.pos])
	l.pos++

	l.pos++ // \n
	l.pos++ // 空格
	l.pos++ // l

	a.False(l.match("2222")) // 不匹配，不会移动位置
	a.True(l.match("ine2"))  // 正确匹配
	l.backup()
	l.backup()
	a.Equal('i', l.data[l.pos])
	l.pos++

	// 超过剩余字符的长度。
	a.False(l.match("ne2\n\n"))
}

func TestLexer_skipSpace(t *testing.T) {
	a := assert.New(t)

	l := newLexer([]rune(" line1\n line2 \n"))
	a.NotNil(l)

	l.skipSpace()
	a.Equal(l.data[l.pos:], "line1\n line2 \n")

	a.True(l.match("line1"))
	a.Equal(l.data[l.pos:], "\n line2 \n")
	l.skipSpace()
	a.Equal(l.data[l.pos:], "line2 \n")
}

func TestLexer_next(t *testing.T) {
	a := assert.New(t)

	l := newLexer([]rune(" \n"))
	a.NotNil(l)

	l.next()
	a.Equal(1, l.pos)

	l.next()
	a.Equal(2, l.pos)

	l.next()
	l.next()
	l.next()
	a.Equal(2, l.pos)
}

func TestTag_readWord(t *testing.T) {
	a := assert.New(t)

	l := &tag{data: []rune(" line1\n line2 \n")}
	a.NotNil(l)

	a.Equal(l.readWord(), []rune("line1"))
	a.Equal(l.readWord(), []rune("line2"))
}

func TestTag_readLine(t *testing.T) {
	a := assert.New(t)

	l := &tag{data: []rune(" line1\n line2 \n")}
	a.NotNil(l)

	a.Equal(l.readLine(), []rune("line1"))
	a.Equal(l.readLine(), []rune("line2"))
}

func TestTag_readEnd(t *testing.T) {
	a := assert.New(t)

	l := &tag{data: []rune(" line1\n line2 \n")}
	a.NotNil(l)

	a.Equal(l.readEnd(), "line1\n line2 \n")
	a.Equal(l.readEnd(), "")
}

func TestTrimRight(t *testing.T) {
	a := assert.New(t)

	a.Equal(trimRight([]rune("123  ")), []rune("123"))
	a.Equal(trimRight([]rune("\n123  ")), []rune("\n123"))
	a.Equal(trimRight([]rune("123\n  ")), []rune("123"))
	a.Equal(trimRight([]rune("123 \n  ")), []rune("123"))
}

// go1.6 BenchmarkTag_readWord-4	10000000	       131 ns/op
func BenchmarkTag_readWord(b *testing.B) {
	a := assert.New(b)
	t := &tag{data: []rune("line1\n @delimiter line2 \n")}
	a.NotNil(t)

	for i := 0; i < b.N; i++ {
		_ = t.readWord()
		t.pos = 0
	}
}

// go1.6 BenchmarkTag_readLine-4	10000000	       109 ns/op
func BenchmarkTag_readLine(b *testing.B) {
	a := assert.New(b)
	t := &tag{data: []rune("line1\n @delimiter line2 \n")}
	a.NotNil(t)

	for i := 0; i < b.N; i++ {
		_ = t.readLine()
		t.pos = 0
	}
}

// go1.6 BenchmarkTag_readEnd-4 	10000000	       181 ns/op
func BenchmarkTag_readEnd(b *testing.B) {
	a := assert.New(b)
	t := &tag{data: []rune("line1\n line2 \n")}
	a.NotNil(t)

	for i := 0; i < b.N; i++ {
		_ = t.readEnd()
		t.pos = 0
	}
}

// go1.6 BenchmarkNewLexer-4    	300000000	         5.63 ns/op
func BenchmarkNewLexer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = newLexer([]rune("line"))
	}
}
